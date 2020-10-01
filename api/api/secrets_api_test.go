package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gojek/mlp/api/models"
	"github.com/gojek/mlp/api/service/mocks"
)

func TestCreateSecret(t *testing.T) {
	testCases := []struct {
		desc               string
		vars               map[string]string
		existingProject    *models.Project
		errFetchingProject error
		body               interface{}
		savedSecret        *models.Secret
		errSaveSecret      error
		expectedResponse   *ApiResponse
	}{
		{
			desc: "Should success",
			vars: map[string]string{
				"project_id": "1",
			},
			body: &models.Secret{
				Name: "name",
				Data: `{"id": 3}`,
			},
			expectedResponse: &ApiResponse{
				code: 201,
				data: &models.Secret{
					Id:        models.Id(1),
					ProjectId: models.Id(1),
					Name:      "name",
					Data:      "encryptedData",
				},
			},
			existingProject: &models.Project{
				Id:   models.Id(1),
				Name: "project",
			},
			savedSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "encryptedData",
			},
		},
		{
			desc: "Should return not found if project is not exist",
			vars: map[string]string{
				"project_id": "1",
			},
			body: &models.Secret{
				Name: "name",
				Data: `{"id": 3}`,
			},
			expectedResponse: &ApiResponse{
				code: 404,
				data: ErrorMessage{"Project with given `project_id: 1` not found"},
			},
			errFetchingProject: fmt.Errorf("project not found"),
		},
		{
			desc: "Should got bad request when body is not complete",
			vars: map[string]string{
				"project_id": "1",
			},
			body: &models.Secret{
				Name: "name",
			},

			expectedResponse: &ApiResponse{
				code: 400,
				data: ErrorMessage{"Invalid request body"},
			},
			existingProject: &models.Project{
				Id:   models.Id(1),
				Name: "project",
			},
		},
		{
			desc: "Should return internal server error when failed save secret",
			vars: map[string]string{
				"project_id": "1",
			},
			body: &models.Secret{
				Name: "name",
				Data: `{"id": 3}`,
			},
			expectedResponse: &ApiResponse{
				code: 500,
				data: ErrorMessage{"db is down"},
			},
			existingProject: &models.Project{
				Id:   models.Id(1),
				Name: "project",
			},
			errSaveSecret: fmt.Errorf("db is down"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			projectService := &mocks.ProjectsService{}
			projectService.On("FindById", models.Id(1)).Return(tC.existingProject, tC.errFetchingProject)

			secretService := &mocks.SecretService{}
			secretService.On("Save", mock.Anything).Return(tC.savedSecret, tC.errSaveSecret)

			controller := &SecretsController{
				AppContext: &AppContext{
					SecretService:   secretService,
					ProjectsService: projectService,
				},
			}

			res := controller.CreateSecret(&http.Request{}, tC.vars, tC.body)
			assert.Equal(t, tC.expectedResponse, res)
		})
	}
}

func TestUpdateSecret(t *testing.T) {
	testCases := []struct {
		desc              string
		vars              map[string]string
		body              interface{}
		existingSecret    *models.Secret
		errFetchingSecret error
		updatedSecret     *models.Secret
		errUpdatingSecret error
		expectedResponse  *ApiResponse
	}{
		{
			desc: "Should responded 204",
			vars: map[string]string{
				"project_id": "1",
				"secret_id":  "1",
			},
			body: &models.Secret{
				Name: "name",
				Data: `{"id": 3}`,
			},
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "no-name",
				Data:      `{"id": 2}`,
			},
			updatedSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      `{"id": 3}`,
			},
			expectedResponse: &ApiResponse{
				code: 200,
				data: &models.Secret{
					Id:        models.Id(1),
					ProjectId: models.Id(1),
					Name:      "name",
					Data:      `{"id": 3}`,
				},
			},
		},
		{
			desc: "Should responded 204 even body is partially there",
			vars: map[string]string{
				"project_id": "1",
				"secret_id":  "1",
			},
			body: &models.Secret{
				Name: "name",
			},
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "no-name",
				Data:      `{"id": 2}`,
			},
			updatedSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      `{"id": 3}`,
			},
			expectedResponse: &ApiResponse{
				code: 200,
				data: &models.Secret{
					Id:        models.Id(1),
					ProjectId: models.Id(1),
					Name:      "name",
					Data:      `{"id": 3}`,
				},
			},
		},
		{
			desc: "Should responded 400 when project_id and secret_id is not integer",
			vars: map[string]string{
				"project_id": "abc",
				"secret_id":  "def",
			},
			body: &models.Secret{
				Name: "name",
			},
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "no-name",
				Data:      `{"id": 2}`,
			},
			updatedSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      `{"id": 3}`,
			},
			expectedResponse: &ApiResponse{
				code: 400,
				data: ErrorMessage{"project_id and secret_id is not valid"},
			},
		},
		{
			desc: "Should responded 400 when body is invalid",
			vars: map[string]string{
				"project_id": "1",
				"secret_id":  "1",
			},
			body: "body",
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "no-name",
				Data:      `{"id": 2}`,
			},
			updatedSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      `{"id": 3}`,
			},
			expectedResponse: &ApiResponse{
				code: 400,
				data: ErrorMessage{"Invalid request body"},
			},
		},
		{
			desc: "Should responded 404 when secret not found",
			vars: map[string]string{
				"project_id": "1",
				"secret_id":  "1",
			},
			body: &models.Secret{
				Name: "name",
			},
			existingSecret:    nil,
			errFetchingSecret: nil,
			expectedResponse: &ApiResponse{
				code: 404,
				data: ErrorMessage{"Secret with given `secret_id: 1` and `project_id: 1` not found"},
			},
		},
		{
			desc: "Should responded 500 when error fetching secret",
			vars: map[string]string{
				"project_id": "1",
				"secret_id":  "1",
			},
			body: &models.Secret{
				Name: "name",
			},
			errFetchingSecret: fmt.Errorf("db is down"),
			expectedResponse: &ApiResponse{
				code: 500,
				data: ErrorMessage{"db is down"},
			},
		},
		{
			desc: "Should responded 500 when error fetching secret",
			vars: map[string]string{
				"project_id": "1",
				"secret_id":  "1",
			},
			body: &models.Secret{
				Name: "name",
			},
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "no-name",
				Data:      `{"id": 2}`,
			},
			errUpdatingSecret: fmt.Errorf("db is down"),
			expectedResponse: &ApiResponse{
				code: 500,
				data: ErrorMessage{"db is down"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			secretService := &mocks.SecretService{}
			secretService.On("FindByIdAndProjectId", models.Id(1), models.Id(1)).Return(tC.existingSecret, tC.errFetchingSecret)
			secretService.On("Save", mock.Anything).Return(tC.updatedSecret, tC.errUpdatingSecret)

			controller := &SecretsController{
				AppContext: &AppContext{
					SecretService: secretService,
				},
			}

			res := controller.UpdateSecret(&http.Request{}, tC.vars, tC.body)
			assert.Equal(t, tC.expectedResponse, res)
		})
	}
}

func TestDeleteSecret(t *testing.T) {
	testCases := []struct {
		desc              string
		vars              map[string]string
		errDeletingSecret error
		expectedResponse  *ApiResponse
	}{
		{
			desc: "Should responsed 204",
			vars: map[string]string{
				"project_id": "1",
				"secret_id":  "1",
			},
			expectedResponse: &ApiResponse{
				code: 204,
				data: nil,
			},
		},
		{
			desc: "Should responsed 400 if project_id or secret_id is invalid",
			vars: map[string]string{
				"project_id": "def",
				"secret_id":  "ghi",
			},
			expectedResponse: &ApiResponse{
				code: 400,
				data: ErrorMessage{"project_id and secret_id is not valid"},
			},
		},
		{
			desc: "Should responsed 500",
			vars: map[string]string{
				"project_id": "1",
				"secret_id":  "1",
			},
			errDeletingSecret: fmt.Errorf("db is down"),
			expectedResponse: &ApiResponse{
				code: 500,
				data: ErrorMessage{"db is down"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			secretService := &mocks.SecretService{}
			secretService.On("Delete", models.Id(1), models.Id(1)).Return(tC.errDeletingSecret)

			controller := &SecretsController{
				AppContext: &AppContext{
					SecretService: secretService,
				},
			}
			res := controller.DeleteSecret(&http.Request{}, tC.vars, nil)
			assert.Equal(t, tC.expectedResponse, res)
		})
	}
}

func TestListSecret(t *testing.T) {
	testCases := []struct {
		desc               string
		vars               map[string]string
		existingProject    *models.Project
		errFetchingProject error
		secrets            []*models.Secret
		errSaveSecret      error
		expectedResponse   *ApiResponse
	}{
		{
			desc: "Should success",
			vars: map[string]string{
				"project_id": "1",
			},
			expectedResponse: &ApiResponse{
				code: 200,
				data: []*models.Secret{
					{
						Id:        models.Id(1),
						ProjectId: models.Id(1),
						Name:      "name-1",
						Data:      "encryptedData",
					},
					{
						Id:        models.Id(2),
						ProjectId: models.Id(1),
						Name:      "name-2",
						Data:      "encryptedData",
					},
				},
			},
			existingProject: &models.Project{
				Id:   models.Id(1),
				Name: "project",
			},
			secrets: []*models.Secret{
				{
					Id:        models.Id(1),
					ProjectId: models.Id(1),
					Name:      "name-1",
					Data:      "encryptedData",
				},
				{
					Id:        models.Id(2),
					ProjectId: models.Id(1),
					Name:      "name-2",
					Data:      "encryptedData",
				},
			},
		},
		{
			desc: "Should return not found if project is not exist",
			vars: map[string]string{
				"project_id": "1",
			},
			expectedResponse: &ApiResponse{
				code: 404,
				data: ErrorMessage{"Project with given `project_id: 1` not found"},
			},
			errFetchingProject: fmt.Errorf("project not found"),
		},
		{
			desc: "Should return internal server error when listing secrets",
			vars: map[string]string{
				"project_id": "1",
			},
			expectedResponse: &ApiResponse{
				code: 500,
				data: ErrorMessage{"db is down"},
			},
			existingProject: &models.Project{
				Id:   models.Id(1),
				Name: "project",
			},
			errSaveSecret: fmt.Errorf("db is down"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			projectService := &mocks.ProjectsService{}
			projectService.On("FindById", models.Id(1)).Return(tC.existingProject, tC.errFetchingProject)

			secretService := &mocks.SecretService{}
			secretService.On("ListSecret", mock.Anything).Return(tC.secrets, tC.errSaveSecret)

			controller := &SecretsController{
				AppContext: &AppContext{
					SecretService:   secretService,
					ProjectsService: projectService,
				},
			}

			res := controller.ListSecret(&http.Request{}, tC.vars, nil)
			assert.Equal(t, tC.expectedResponse, res)
		})
	}
}
