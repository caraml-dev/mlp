package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"

	"github.com/caraml-dev/mlp/api/models"
)

func (s *APITestSuite) TestCreateSecretStorage() {
	type args struct {
		path string
		body interface{}
	}

	nonExistingProjectID := models.ID(123)

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "success: create project-scoped secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages", s.mainProject.ID),
				body: &models.SecretStorage{
					Name:      "new-secret-storage",
					Type:      models.VaultSecretStorageType,
					Scope:     models.ProjectSecretStorageScope,
					ProjectID: &s.mainProject.ID,
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:        "http://localhost:8200",
							Role:       "my-role",
							MountPath:  "secret",
							PathPrefix: fmt.Sprintf("secret-storage/%d", time.Now().Unix()),
							AuthMethod: models.TokenAuthMethod,
							Token:      "root",
						},
					},
				},
			},
			want: &Response{
				code: http.StatusCreated,
				data: &models.SecretStorage{
					Name:      "new-secret-storage",
					Type:      models.VaultSecretStorageType,
					Scope:     models.ProjectSecretStorageScope,
					ProjectID: &s.mainProject.ID,
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:        "http://localhost:8200",
							Role:       "my-role",
							MountPath:  "secret",
							PathPrefix: fmt.Sprintf("secret-storage/%d", time.Now().Unix()),
							AuthMethod: models.TokenAuthMethod,
							Token:      "root",
						},
					},
				},
			},
		},
		{
			name: "error: create internal secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages", s.mainProject.ID),
				body: &models.SecretStorage{
					Name:      "new-secret-storage",
					Type:      models.InternalSecretStorageType,
					Scope:     models.ProjectSecretStorageScope,
					ProjectID: &s.mainProject.ID,
				},
			},
			want: &Response{
				code: http.StatusBadRequest,
				data: ErrorMessage{
					Message: "invalid secret storage type: internal",
				},
			},
		},
		{
			name: "error: create global secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages", s.mainProject.ID),
				body: &models.SecretStorage{
					Name:      "new-secret-storage",
					Type:      models.VaultSecretStorageType,
					Scope:     models.GlobalSecretStorageScope,
					ProjectID: &s.mainProject.ID,
				},
			},
			want: &Response{
				code: http.StatusBadRequest,
				data: ErrorMessage{
					Message: "invalid secret storage scope: global",
				},
			},
		},
		{
			name: "error: create secret storage in non existing project",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages", nonExistingProjectID),
				body: &models.SecretStorage{
					Name:      "new-secret-storage",
					Type:      models.VaultSecretStorageType,
					Scope:     models.ProjectSecretStorageScope,
					ProjectID: &nonExistingProjectID,
				},
			},
			want: &Response{
				code: http.StatusNotFound,
				data: ErrorMessage{
					Message: "project with ID 123 not found",
				},
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
				var ss models.SecretStorage
				jsonObj.Decode(&ss)
				assertSecretStorageEquals(s.T(), tt.want.data.(*models.SecretStorage), &ss)
			} else {
				var err ErrorMessage
				jsonObj.Decode(&err)
				s.Equal(tt.want.data, err)
			}
		})
	}
}

func (s *APITestSuite) TestListSecretStorage() {
	type args struct {
		path string
	}

	nonExistingProjectID := models.ID(123)

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "success: list all secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages", s.mainProject.ID),
			},
			want: &Response{
				code: http.StatusOK,
				data: []*models.SecretStorage{
					s.internalSecretStorage, s.defaultSecretStorage, s.projectSecretStorage,
				},
			},
		},
		{
			name: "error: list secret storage in non existing project",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages", nonExistingProjectID),
			},
			want: &Response{
				code: http.StatusNotFound,
				data: ErrorMessage{
					Message: "project with ID 123 not found",
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			server := httptest.NewServer(s.route)
			defer server.Close()

			e := httpexpect.Default(s.T(), server.URL)
			j := e.GET(tt.args.path).
				Expect().
				Status(tt.want.code).
				JSON()

			if tt.want.code >= http.StatusOK && tt.want.code < http.StatusMultipleChoices {
				var secretStorages []*models.SecretStorage
				j.Array().Decode(&secretStorages)
				for i, ss := range tt.want.data.([]*models.SecretStorage) {
					assertSecretStorageEquals(s.T(), ss, secretStorages[i])
				}
			} else {
				var err ErrorMessage
				j.Object().Decode(&err)
				s.Equal(tt.want.data, err)
			}
		})
	}
}

func (s *APITestSuite) TestGetSecretStorage() {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "success: get project-scoped secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.projectSecretStorage.ID),
			},
			want: &Response{
				code: http.StatusOK,
				data: s.projectSecretStorage,
			},
		},
		{
			name: "success: get default secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.defaultSecretStorage.ID),
			},
			want: &Response{
				code: http.StatusOK,
				data: s.defaultSecretStorage,
			},
		},
		{
			name: "success: get internal secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.internalSecretStorage.ID),
			},
			want: &Response{
				code: http.StatusOK,
				data: s.internalSecretStorage,
			},
		},
		{
			name: "error: get non existing secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, 123),
			},
			want: &Response{
				code: http.StatusNotFound,
				data: ErrorMessage{
					Message: "secret storage with ID 123 not found",
				},
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
				var secretStorage *models.SecretStorage
				jsonObj.Decode(&secretStorage)
				assertSecretStorageEquals(s.T(), tt.want.data.(*models.SecretStorage), secretStorage)
			} else {
				var err ErrorMessage
				jsonObj.Decode(&err)
				s.Equal(tt.want.data, err)
			}
		})
	}
}

func (s *APITestSuite) TestDeleteSecretStorage() {
	// 1. success: delete project-scoped secret storage
	// 2. success: delete non existing secret storage
	// 3. success: delete internal secret storage
	// 4. failure: delete default secret storage

	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "success: delete project-scoped secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.projectSecretStorage.ID),
			},
			want: &Response{
				code: http.StatusNoContent,
			},
		},
		{
			name: "success: delete non existing secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, 123),
			},
			want: &Response{
				code: http.StatusNoContent,
			},
		},
		{
			name: "success: delete internal secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.internalSecretStorage.ID),
			},
			want: &Response{
				code: http.StatusNoContent,
			},
		},
		{
			name: "error: delete default secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.defaultSecretStorage.ID),
			},
			want: &Response{
				code: http.StatusBadRequest,
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

func (s *APITestSuite) TestUpdateSecretStorage() {
	// 1. success: update project-scoped secret storage name
	// 2. success: migrate project-scoped secret storage to other secret storage
	// 3. error: update non-existing secret storage
	// 4. error: update internal secret storage
	// 5. error: update global secret storage

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
			name: "success: update project-scoped secret storage name",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.projectSecretStorage.ID),
				body: &models.SecretStorage{
					Name: "updated-name",
				},
			},
			want: &Response{
				code: http.StatusOK,
				data: &models.SecretStorage{
					ID:        s.projectSecretStorage.ID,
					Name:      "updated-name",
					ProjectID: &s.mainProject.ID,
					Type:      s.projectSecretStorage.Type,
					Scope:     s.projectSecretStorage.Scope,
					Config:    s.projectSecretStorage.Config,
				},
			},
		},
		{
			name: "success: migrate project-scoped secret storage to other secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.projectSecretStorage.ID),
				body: &models.SecretStorage{
					Name: "project-secret-storage",
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:        "http://localhost:8200",
							Role:       "my-role",
							MountPath:  "secret",
							PathPrefix: "new-mount-path",
							AuthMethod: models.TokenAuthMethod,
							Token:      "root",
						},
					},
				},
			},
			want: &Response{
				code: http.StatusOK,
				data: &models.SecretStorage{
					ID:        s.projectSecretStorage.ID,
					Name:      "project-secret-storage",
					ProjectID: &s.mainProject.ID,
					Type:      s.projectSecretStorage.Type,
					Scope:     s.projectSecretStorage.Scope,
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:        "http://localhost:8200",
							Role:       "my-role",
							MountPath:  "secret",
							PathPrefix: "new-mount-path",
							AuthMethod: models.TokenAuthMethod,
							Token:      "root",
						},
					},
				},
			},
		},
		{
			name: "error: update non-existing secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, 123),
				body: &models.SecretStorage{
					Name: "updated-name",
				},
			},
			want: &Response{
				code: http.StatusNotFound,
				data: ErrorMessage{
					Message: "secret storage with ID 123 not found",
				},
			},
		},
		{
			name: "error: update global secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.defaultSecretStorage.ID),
				body: &models.SecretStorage{
					Name: "updated-name",
				},
			},
			want: &Response{
				code: http.StatusBadRequest,
				data: ErrorMessage{
					Message: "cannot update global secret storage",
				},
			},
		},
		{
			name: "error: update internal secret storage",
			args: args{
				path: fmt.Sprintf("/v1/projects/%d/secret_storages/%d", s.mainProject.ID, s.internalSecretStorage.ID),
				body: &models.SecretStorage{
					Name: "updated-name",
				},
			},
			want: &Response{
				code: http.StatusBadRequest,
				data: ErrorMessage{
					Message: "cannot update global secret storage",
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
				var secretStorage *models.SecretStorage
				jsonObj.Decode(&secretStorage)
				assertSecretStorageEquals(s.T(), tt.want.data.(*models.SecretStorage), secretStorage)
			} else {
				var err ErrorMessage
				jsonObj.Decode(&err)
				s.Equal(tt.want.data, err)
			}
		})
	}
}

func assertSecretStorageEquals(t *testing.T, want *models.SecretStorage, got *models.SecretStorage) {
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Type, got.Type)
	assert.Equal(t, want.Scope, got.Scope)
	assert.Equal(t, want.ProjectID, got.ProjectID)
	assert.Equal(t, want.Config, got.Config)

	assert.NotEmpty(t, got.ID)
	assert.NotEmpty(t, got.CreatedAt)
	assert.NotEmpty(t, got.UpdatedAt)
}
