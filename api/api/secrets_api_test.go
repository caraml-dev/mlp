//go:build integration

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/it/database"
	"github.com/caraml-dev/mlp/api/models"
)

type SecretAPITestSuite struct {
	suite.Suite
	cleanupFn func()
	route     http.Handler

	internalSecretStorage *models.SecretStorage
	defaultSecretStorage  *models.SecretStorage
	project               *models.Project
	existingSecrets       []*models.Secret
}

func (s *SecretAPITestSuite) SetupTest() {
	db, cleanupFn, err := database.CreateTestDatabase()
	s.Require().NoError(err, "Failed to connect to test database")
	s.cleanupFn = cleanupFn

	appCtx, err := NewAppContext(db, &config.Config{
		Port: 0,
		Authorization: &config.AuthorizationConfig{
			Enabled: false,
		},
		Mlflow: &config.MlflowConfig{
			TrackingURL: "http://mlflow:5000",
		},
		DefaultSecretStorage: &config.SecretStorage{
			Name: "vault",
			Type: string(models.VaultSecretStorageType),
			Config: models.SecretStorageConfig{
				VaultConfig: &models.VaultConfig{
					URL:        "http://localhost:8200",
					Role:       "my-role",
					MountPath:  "secret",
					PathPrefix: fmt.Sprintf("secret-api-test/%d/{{ .Project }}", time.Now().Unix()),
					AuthMethod: models.TokenAuthMethod,
					Token:      "root",
				},
			},
		},
	})
	s.Require().NoError(err, "Failed to create app context")

	s.internalSecretStorage, err = appCtx.SecretStorageService.FindByID(1)
	s.Require().NoError(err, "Failed to find internal secret storage")

	s.defaultSecretStorage = appCtx.DefaultSecretStorage

	s.project, err = appCtx.ProjectsService.CreateProject(&models.Project{
		Name: "test-project",
	})
	s.Require().NoError(err, "Failed to create project")

	controller := &SecretsController{
		AppContext: appCtx,
	}

	s.existingSecrets = make([]*models.Secret, 0)
	for i := 0; i < 5; i++ {
		secretStorageID := s.defaultSecretStorage.ID
		if i < 2 {
			secretStorageID = s.internalSecretStorage.ID
		}

		secret, err := appCtx.SecretService.Create(&models.Secret{
			ProjectID:       s.project.ID,
			SecretStorageID: &secretStorageID,
			Name:            fmt.Sprintf("secret-%d", i),
			Data:            fmt.Sprintf("secret-data-%d", i),
		})
		s.Require().NoError(err, "Failed to create secret")
		s.existingSecrets = append(s.existingSecrets, secret)
	}

	controllers := []Controller{controller}
	r := NewRouter(appCtx, controllers)

	route := mux.NewRouter()
	route.PathPrefix(basePath).Handler(
		http.StripPrefix(
			strings.TrimSuffix(basePath, "/"),
			r,
		),
	)

	s.route = route
}

func (s *SecretAPITestSuite) TearDownTest() {
	s.cleanupFn()
}

func (s *SecretAPITestSuite) TestGetSecret() {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "success: get internal secret",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.project.ID, s.existingSecrets[0].ID),
			},
			want: &Response{
				code: http.StatusOK,
				data: s.existingSecrets[0],
			},
		},
		{
			name: "success: get external secret",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.project.ID, s.existingSecrets[2].ID),
			},
			want: &Response{
				code: http.StatusOK,
				data: s.existingSecrets[2],
			},
		},
		{
			name: "error: secret not found",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", 123, 123),
			},
			want: &Response{
				code: http.StatusNotFound,
				data: ErrorMessage{"Secret with given `secret_id: 123` not found"},
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {

			server := httptest.NewServer(s.route)
			defer server.Close()

			e := httpexpect.Default(s.T(), server.URL)

			jsonObj := e.GET(tt.args.path).
				Expect().
				Status(tt.want.code).
				JSON().Object()

			if tt.want.code >= http.StatusOK && tt.want.code < http.StatusMultipleChoices {
				var secret models.Secret
				jsonObj.Decode(&secret)
				assertSecretEquals(s.T(), tt.want.data.(*models.Secret), &secret)
			} else {
				var err ErrorMessage
				jsonObj.Decode(&err)
				s.Equal(tt.want.data, err)
			}
		})
	}
}

func (s *SecretAPITestSuite) TestCreateSecret() {
	type args struct {
		path string
		body interface{}
	}

	nonExistingSecretStorageID := models.ID(123)

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "success: create secret in default secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.project.ID),
				body: &models.Secret{
					Name: "external-secrete",
					Data: "secret-value",
				},
			},
			want: &Response{
				code: http.StatusCreated,
				data: &models.Secret{
					ID:              models.ID(6),
					SecretStorageID: &s.defaultSecretStorage.ID,
					ProjectID:       s.project.ID,
					Name:            "external-secrete",
					Data:            "secret-value",
				},
			},
		},
		{
			name: "success: create secret in internal secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.project.ID),
				body: &models.Secret{
					Name:            "internal-secret",
					Data:            "secret-value",
					SecretStorageID: &s.internalSecretStorage.ID,
				},
			},
			want: &Response{
				code: http.StatusCreated,
				data: &models.Secret{
					ID:              models.ID(7),
					SecretStorageID: &s.internalSecretStorage.ID,
					ProjectID:       s.project.ID,
					Name:            "internal-secret",
					Data:            "secret-value",
				},
			},
		},
		{
			name: "error: should return not found if project is not exist",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", 123),
				body: &models.Secret{
					Name: "internal-secret",
					Data: "secret-value",
				},
			},
			want: &Response{
				code: http.StatusNotFound,
				data: ErrorMessage{"Project with given `project_id: 123` not found"},
			},
		},
		{
			name: "error: should return not found if secret storage is not exist",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.project.ID),
				body: &models.Secret{
					Name:            "internal-secret",
					Data:            "secret-value",
					SecretStorageID: &nonExistingSecretStorageID,
				},
			},
			want: &Response{
				code: http.StatusNotFound,
				data: ErrorMessage{"Secret storage with given `secret_storage_id: 123` not found"},
			},
		},
		{
			name: "error: should got bad request when body is not complete",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.project.ID),
				body: &models.Secret{
					Name: "internal-secret",
				},
			},
			want: &Response{
				code: http.StatusBadRequest,
				data: ErrorMessage{"Invalid request body"},
			},
		},
		{
			name: "error: should return secret already exist for same name",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.project.ID),
				body: &models.Secret{
					Name: "external-secrete",
					Data: "secret-value",
				},
			},
			want: &Response{
				code: http.StatusInternalServerError,
				data: ErrorMessage{"error when saving secret in database, " +
					"error: pq: duplicate key value violates unique constraint \"secrets_project_id_name_key\""},
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {

			server := httptest.NewServer(s.route)
			defer server.Close()

			e := httpexpect.Default(s.T(), server.URL)

			jsonObj := e.POST(tt.args.path).
				WithJSON(tt.args.body).
				Expect().
				Status(tt.want.code).
				JSON().Object()

			if tt.want.code >= http.StatusOK && tt.want.code < http.StatusMultipleChoices {
				var secret models.Secret
				jsonObj.Decode(&secret)
				assertSecretEquals(s.T(), tt.want.data.(*models.Secret), &secret)
			} else {
				var err ErrorMessage
				jsonObj.Decode(&err)
				s.Equal(tt.want.data, err)
			}
		})
	}
}

func (s *SecretAPITestSuite) TestUpdateSecret() {
	type args struct {
		path string
		body interface{}
	}

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "success: update secret value",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.project.ID, s.existingSecrets[2].ID),
				body: &models.Secret{
					Name: s.existingSecrets[2].Name,
					Data: "new-value",
				},
			},
			want: &Response{
				code: http.StatusOK,
				data: &models.Secret{
					ID:              s.existingSecrets[2].ID,
					SecretStorageID: s.existingSecrets[2].SecretStorageID,
					ProjectID:       s.project.ID,
					Name:            s.existingSecrets[2].Name,
					Data:            "new-value",
				},
			},
		},
		{
			name: "success: migrate secret",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.project.ID, s.existingSecrets[0].ID),
				body: &models.Secret{
					SecretStorageID: &s.defaultSecretStorage.ID,
				},
			},
			want: &Response{
				code: http.StatusOK,
				data: &models.Secret{
					ID:              s.existingSecrets[0].ID,
					SecretStorageID: &s.defaultSecretStorage.ID,
					ProjectID:       s.project.ID,
					Name:            s.existingSecrets[0].Name,
					Data:            s.existingSecrets[0].Data,
				},
			},
		},
		{
			name: "success: migrate secret and update value",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.project.ID, s.existingSecrets[1].ID),
				body: &models.Secret{
					Name:            "secret-1",
					Data:            "new-value",
					SecretStorageID: &s.defaultSecretStorage.ID,
				},
			},
			want: &Response{
				code: http.StatusOK,
				data: &models.Secret{
					ID:              s.existingSecrets[1].ID,
					SecretStorageID: &s.defaultSecretStorage.ID,
					ProjectID:       s.project.ID,
					Name:            "secret-1",
					Data:            "new-value",
				},
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {

			server := httptest.NewServer(s.route)
			defer server.Close()

			e := httpexpect.Default(s.T(), server.URL)

			jsonObj := e.PATCH(tt.args.path).
				WithJSON(tt.args.body).
				Expect().
				Status(tt.want.code).
				JSON().Object()

			if tt.want.code >= http.StatusOK && tt.want.code < http.StatusMultipleChoices {
				var secret models.Secret
				jsonObj.Decode(&secret)
				assertSecretEquals(s.T(), tt.want.data.(*models.Secret), &secret)
			} else {
				var err ErrorMessage
				jsonObj.Decode(&err)
				s.Equal(tt.want.data, err)
			}
		})
	}
}

func (s *SecretAPITestSuite) TestDeleteSecret() {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "success: delete internal secret",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.project.ID, s.existingSecrets[0].ID),
			},
			want: &Response{
				code: http.StatusNoContent,
			},
		},
		{
			name: "success: delete external secret",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.project.ID, s.existingSecrets[2].ID),
			},
			want: &Response{
				code: http.StatusNoContent,
			},
		},
		{
			name: "success: secret not found",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", 123, 123),
			},
			want: &Response{
				code: http.StatusNoContent,
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {

			server := httptest.NewServer(s.route)
			defer server.Close()

			e := httpexpect.Default(s.T(), server.URL)

			e.DELETE(tt.args.path).
				Expect().
				Status(tt.want.code)
		})
	}
}

func (s *SecretAPITestSuite) TestListSecret() {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "success: get internal secret",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.project.ID),
			},
			want: &Response{
				code: http.StatusOK,
				data: s.existingSecrets,
			},
		},
		{
			name: "error: project not found",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", 123),
			},
			want: &Response{
				code: http.StatusNotFound,
				data: ErrorMessage{"Project with given `project_id: 123` not found"},
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {

			server := httptest.NewServer(s.route)
			defer server.Close()

			e := httpexpect.Default(s.T(), server.URL)

			jsonObj := e.GET(tt.args.path).
				Expect().
				Status(tt.want.code).
				JSON()

			if tt.want.code >= http.StatusOK && tt.want.code < http.StatusMultipleChoices {
				var secrets []*models.Secret
				jsonObj.Array().Decode(&secrets)
				for i, secret := range tt.want.data.([]*models.Secret) {
					assertSecretEquals(s.T(), secret, secrets[i])
				}
			} else {
				var err ErrorMessage
				jsonObj.Object().Decode(&err)
				s.Equal(tt.want.data, err)
			}
		})
	}
}

func TestSecretAPI(t *testing.T) {
	suite.Run(t, new(SecretAPITestSuite))
}

func assertSecretEquals(t *testing.T, exp *models.Secret, got *models.Secret) {
	assert.Equal(t, exp.ID, got.ID)
	assert.Equal(t, exp.ProjectID, got.ProjectID)
	assert.Equal(t, exp.Name, got.Name)
	assert.Equal(t, exp.Data, got.Data)
	assert.Equal(t, *exp.SecretStorageID, *got.SecretStorageID)

	assert.NotEmpty(t, got.CreatedAt)
	assert.NotEmpty(t, got.UpdatedAt)
}
