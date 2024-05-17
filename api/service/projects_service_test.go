package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"

	"github.com/stretchr/testify/assert"

	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/webhooks"
	"github.com/caraml-dev/mlp/api/repository/mocks"

	enforcerMock "github.com/caraml-dev/mlp/api/pkg/authz/enforcer/mocks"
)

const MLFlowTrackingURL = "http://localhost:5555"

func TestProjectsService_CreateProject(t *testing.T) {
	tests := []struct {
		name          string
		arg           *models.Project
		authEnabled   bool
		expResult     *models.Project
		expAuthUpdate *enforcer.AuthorizationUpdateRequest
		wantError     bool
		wantErrorMsg  string
	}{
		{
			"success: auth enabled",
			&models.Project{
				ID:             1,
				Name:           "my-project",
				Administrators: []string{"user@email.com"},
				Readers:        nil,
			},
			true,
			&models.Project{
				ID:                1,
				Name:              "my-project",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Administrators:    []string{"user@email.com"},
				Readers:           nil,
			},
			&enforcer.AuthorizationUpdateRequest{
				RolePermissions: map[string][]string{
					"mlp.administrator": {"mlp.projects.1.get", "mlp.projects.1.put", "mlp.projects.1.post",
						"mlp.projects.1.patch", "mlp.projects.1.delete"},
					"mlp.projects.reader":   {"mlp.projects.1.get"},
					"mlp.projects.1.reader": {"mlp.projects.1.get"},
					"mlp.projects.1.administrator": {"mlp.projects.1.get", "mlp.projects.1.put", "mlp.projects.1.post",
						"mlp.projects.1.patch", "mlp.projects.1.delete"},
				},
				RoleMembers: map[string][]string{
					"mlp.projects.1.reader":        {},
					"mlp.projects.1.administrator": {"user@email.com"},
				},
			},
			false,
			"",
		},
		{
			"success: auth disabled",
			&models.Project{
				ID:             1,
				Name:           "my-project",
				Administrators: []string{"user@email.com"},
				Readers:        nil,
			},
			false,
			&models.Project{
				ID:                1,
				Name:              "my-project",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Administrators:    []string{"user@email.com"},
				Readers:           nil,
			},
			nil,
			false,
			"",
		},
		{
			"failed: reserved name",
			&models.Project{
				ID:             1,
				Name:           "infrastructure",
				Administrators: []string{"user@email.com"},
				Readers:        nil,
			},
			false,
			nil,
			nil,
			true,
			"unable to use reserved project name: infrastructure",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mocks.ProjectRepository{}
			storage.On("Save", tt.expResult).Return(tt.expResult, nil)

			authEnforcer := &enforcerMock.Enforcer{}
			projectsService, err := NewProjectsService(MLFlowTrackingURL, storage, authEnforcer, tt.authEnabled, nil)
			require.NoError(t, err)

			if tt.expAuthUpdate != nil {
				authEnforcer.On("UpdateAuthorization", mock.Anything, *tt.expAuthUpdate).Return(nil)
			}

			res, err := projectsService.CreateProject(context.Background(), tt.arg)
			if tt.wantError {
				require.Error(t, err)
				require.Equal(t, tt.wantErrorMsg, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expResult, res)

			storage.AssertExpectations(t)

			if tt.expAuthUpdate != nil {
				authEnforcer.AssertExpectations(t)
			}
		})
	}
}

func TestProjectsService_UpdateProject(t *testing.T) {
	tests := []struct {
		name          string
		arg           *models.Project
		authEnabled   bool
		expResult     *models.Project
		expAuthUpdate *enforcer.AuthorizationUpdateRequest
	}{
		{
			"success: auth enabled",
			&models.Project{
				ID:             1,
				Name:           "my-project",
				Administrators: []string{"user@email.com"},
				Readers:        nil,
			},
			true,
			&models.Project{
				ID:                1,
				Name:              "my-project",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Administrators:    []string{"user@email.com"},
				Readers:           nil,
			},
			&enforcer.AuthorizationUpdateRequest{
				RolePermissions: map[string][]string{
					"mlp.administrator": {"mlp.projects.1.get", "mlp.projects.1.put", "mlp.projects.1.post",
						"mlp.projects.1.patch", "mlp.projects.1.delete"},
					"mlp.projects.reader":   {"mlp.projects.1.get"},
					"mlp.projects.1.reader": {"mlp.projects.1.get"},
					"mlp.projects.1.administrator": {"mlp.projects.1.get", "mlp.projects.1.put", "mlp.projects.1.post",
						"mlp.projects.1.patch", "mlp.projects.1.delete"},
				},
				RoleMembers: map[string][]string{
					"mlp.projects.1.reader":        {},
					"mlp.projects.1.administrator": {"user@email.com"},
				},
			},
		},
		{
			"success: auth disabled",
			&models.Project{
				ID:             1,
				Name:           "my-project",
				Administrators: []string{"user@email.com"},
				Readers:        nil,
			},
			false,
			&models.Project{
				ID:                1,
				Name:              "my-project",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Administrators:    []string{"user@email.com"},
				Readers:           nil,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mocks.ProjectRepository{}
			storage.On("Save", tt.expResult).Return(tt.expResult, nil)

			authEnforcer := &enforcerMock.Enforcer{}
			if tt.expAuthUpdate != nil {
				authEnforcer.On("UpdateAuthorization", mock.Anything, *tt.expAuthUpdate).Return(nil)
			}

			projectsService, err := NewProjectsService(MLFlowTrackingURL, storage, authEnforcer, tt.authEnabled, nil)
			assert.NoError(t, err)

			res, err := projectsService.UpdateProject(context.Background(), tt.arg)
			require.NoError(t, err)
			require.Equal(t, tt.expResult, res)

			storage.AssertExpectations(t)
			authEnforcer.AssertExpectations(t)
		})
	}
}

func TestProjectsService_ListProjects(t *testing.T) {
	project1 := &models.Project{
		ID:                1,
		Name:              "project-1",
		MLFlowTrackingURL: MLFlowTrackingURL,
		Administrators:    []string{"admin-1@email.com"},
		Readers:           []string{"reader-1@email.com"},
	}
	project2 := &models.Project{
		ID:                2,
		Name:              "project-2",
		MLFlowTrackingURL: MLFlowTrackingURL,
		Administrators:    []string{"admin-2@email.com"},
		Readers:           []string{"reader-2@email.com"},
	}
	allProjects := []*models.Project{project1, project2}
	storage := &mocks.ProjectRepository{}
	storage.On("ListProjects", "project-").Return(allProjects, nil)

	tests := []struct {
		name          string
		projectFilter string
		authEnabled   bool
		expResult     []*models.Project
		user          string
		userRoles     []string
	}{
		{
			"filter only by project name when auth is disabled",
			"project-",
			false,
			allProjects,
			"anonymous@email.com",
			nil,
		},
		{
			"filter by permission and project name when auth is enabled",
			"project-",
			true,
			[]*models.Project{},
			"anonymous-user@email.com",
			[]string{},
		},
		{
			"allow project admin to read project, regardless of user roles return by enforcer",
			"project-",
			true,
			[]*models.Project{project1},
			"admin-1@email.com",
			[]string{"some roles"},
		},
		{
			"allow project reader to read project, regardless of user roles return by enforcer",
			"project-",
			true,
			[]*models.Project{project2},
			"reader-2@email.com",
			[]string{"some roles"},
		},
		{
			"allow mlp administrators to read all projects",
			"project-",
			true,
			allProjects,
			"mlp-admin@email.com",
			[]string{"mlp.administrator"},
		},
		{
			"allow project readers to read all projects",
			"project-",
			true,
			allProjects,
			"project-reader@email.com",
			[]string{"mlp.projects.reader"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authEnforcer := &enforcerMock.Enforcer{}
			if tt.authEnabled {
				authEnforcer.On("GetUserRoles", mock.Anything, tt.user).Return(tt.userRoles, nil)
			}

			projectsService, err := NewProjectsService(MLFlowTrackingURL, storage, authEnforcer, tt.authEnabled, nil)
			assert.NoError(t, err)

			res, err := projectsService.ListProjects(context.Background(), "project-", tt.user)
			assert.NoError(t, err)
			assert.Equal(t, tt.expResult, res)

			storage.AssertExpectations(t)
			authEnforcer.AssertExpectations(t)
		})
	}
}

func TestProjectsService_FindById(t *testing.T) {
	id := models.ID(1)

	exp := &models.Project{
		ID:                id,
		Name:              "my-project",
		MLFlowTrackingURL: MLFlowTrackingURL,
		Administrators:    []string{"user@email.com"},
		Readers:           nil,
	}
	storage := &mocks.ProjectRepository{}
	storage.On("Get", id).Return(exp, nil)

	authEnforcer := &enforcerMock.Enforcer{}
	projectsService, err := NewProjectsService(MLFlowTrackingURL, storage, authEnforcer, false, nil)
	assert.NoError(t, err)

	res, err := projectsService.FindByID(id)
	assert.NoError(t, err)
	assert.Equal(t, exp, res)

	storage.AssertExpectations(t)
}

func TestProjectsService_CreateWithWebhook(t *testing.T) {

	authEnforcer := &enforcerMock.Enforcer{}
	mockClient1 := &webhooks.MockWebhookClient{}
	mockClient1.On("IsAsync").Return(false)
	mockClient1.On("GetName").Return("webhook1")
	mockClient1.On("IsFinalResponse").Return(true)
	mockClient1.On("GetUseDataFrom").Return("")
	storage := &mocks.ProjectRepository{}
	tests := []struct {
		name         string
		arg          *models.Project
		authEnabled  bool
		expResult    *models.Project
		wantError    bool
		wantErrorMsg string
		whResponse   []byte
	}{

		{
			name: "test basic working webhook",
			arg: &models.Project{
				ID:                1,
				Name:              "project-1",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Team:              "team-1",
				Stream:            "team-2",
			},
			authEnabled: false,
			expResult: &models.Project{
				ID:                1,
				Name:              "project-1",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Team:              "team-1",
				Stream:            "team-2-modified-by-webhook",
			},
			wantError: false,
			whResponse: []byte(`{
				"id": 1,
				"name": "project-1",
				"team": "team-1",
				"stream": "team-2-modified-by-webhook"
			}`),
		},
		{
			name: "test invalid json str error",
			arg: &models.Project{
				ID:                1,
				Name:              "project-1",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Team:              "team-1",
			},
			authEnabled: false,
			expResult: &models.Project{
				ID:                1,
				Name:              "project-1",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Administrators:    []string{"admin-1@email.com"},
				Readers:           []string{"reader-1@email.com"},
				Team:              "team-1",
				Stream:            "team-2-modified-by-webhook",
			},
			wantError: true,
			whResponse: []byte(`{
				invalid-json-str
			`),
		},
	}

	// construct test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.On("Save", test.arg).Return(test.arg, nil).Once()
			storage.On("Save", test.expResult).Return(test.expResult, nil).Once()
			mockClient1.On("Invoke", mock.Anything, mock.Anything).Return(test.whResponse, nil)
			whManager := &webhooks.SimpleWebhookManager{
				WebhookClients: map[webhooks.EventType][]webhooks.WebhookClient{
					ProjectCreatedEvent: {
						mockClient1,
					},
				},
			}
			mockClient1.On("Invoke", mock.Anything, mock.Anything).Return(test.whResponse, nil)
			projectsService, err := NewProjectsService(MLFlowTrackingURL, storage, authEnforcer, false, whManager)
			assert.NoError(t, err)
			res, err := projectsService.CreateProject(context.Background(), test.arg)
			if test.wantError {
				require.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expResult, res)
		})
	}

}
