//go:build integration

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"

	"github.com/caraml-dev/mlp/api/models"
)

func (s *APITestSuite) TestGetSecret() {
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
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.mainProject.ID, s.existingSecrets[0].ID),
			},
			want: &Response{
				code: http.StatusOK,
				data: s.existingSecrets[0],
			},
		},
		{
			name: "success: get external secret",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.mainProject.ID, s.existingSecrets[2].ID),
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
				data: ErrorMessage{"secret with id 123 not found"},
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

func (s *APITestSuite) TestCreateSecret() {
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
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.mainProject.ID),
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
					ProjectID:       s.mainProject.ID,
					Name:            "external-secrete",
					Data:            "secret-value",
				},
			},
		},
		{
			name: "success: create secret in internal secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.mainProject.ID),
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
					ProjectID:       s.mainProject.ID,
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
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.mainProject.ID),
				body: &models.Secret{
					Name:            "internal-secret",
					Data:            "secret-value",
					SecretStorageID: &nonExistingSecretStorageID,
				},
			},
			want: &Response{
				code: http.StatusNotFound,
				data: ErrorMessage{"error when fetching secret storage with id: 123, error: secret storage with ID 123 not found"},
			},
		},
		{
			name: "error: should got bad request when body is not complete",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.mainProject.ID),
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
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.mainProject.ID),
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

func (s *APITestSuite) TestUpdateSecret() {
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
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.mainProject.ID, s.existingSecrets[2].ID),
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
					ProjectID:       s.mainProject.ID,
					Name:            s.existingSecrets[2].Name,
					Data:            "new-value",
				},
			},
		},
		{
			name: "success: migrate secret",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.mainProject.ID, s.existingSecrets[0].ID),
				body: &models.Secret{
					SecretStorageID: &s.defaultSecretStorage.ID,
				},
			},
			want: &Response{
				code: http.StatusOK,
				data: &models.Secret{
					ID:              s.existingSecrets[0].ID,
					SecretStorageID: &s.defaultSecretStorage.ID,
					ProjectID:       s.mainProject.ID,
					Name:            s.existingSecrets[0].Name,
					Data:            s.existingSecrets[0].Data,
				},
			},
		},
		{
			name: "success: migrate secret and update value",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.mainProject.ID, s.existingSecrets[1].ID),
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
					ProjectID:       s.mainProject.ID,
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

func (s *APITestSuite) TestDeleteSecret() {
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
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.mainProject.ID, s.existingSecrets[0].ID),
			},
			want: &Response{
				code: http.StatusNoContent,
			},
		},
		{
			name: "success: delete external secret",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secrets/%d", s.mainProject.ID, s.existingSecrets[2].ID),
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

func (s *APITestSuite) TestListSecret() {
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
				path: fmt.Sprintf("/v1/projects/%d/secrets", s.mainProject.ID),
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
				data: ErrorMessage{"project with ID 123 not found"},
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

func assertSecretEquals(t *testing.T, exp *models.Secret, got *models.Secret) {
	assert.Equal(t, exp.ID, got.ID)
	assert.Equal(t, exp.ProjectID, got.ProjectID)
	assert.Equal(t, exp.Name, got.Name)
	assert.Equal(t, exp.Data, got.Data)
	assert.Equal(t, *exp.SecretStorageID, *got.SecretStorageID)

	assert.NotEmpty(t, got.CreatedAt)
	assert.NotEmpty(t, got.UpdatedAt)
}
