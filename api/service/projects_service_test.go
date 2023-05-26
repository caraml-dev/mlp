package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
	enforcerMock "github.com/caraml-dev/mlp/api/pkg/authz/enforcer/mocks"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer/types"

	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/repository/mocks"
)

const MLFlowTrackingURL = "http://localhost:5555"

func TestProjectsService_CreateProject(t *testing.T) {
	tests := []struct {
		name         string
		arg          *models.Project
		authEnabled  bool
		expResult    *models.Project
		wantError    bool
		wantErrorMsg string
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
			true,
			"unable to use reserved project name: infrastructure",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mocks.ProjectRepository{}
			storage.On("Save", tt.expResult).Return(tt.expResult, nil)

			projectResource := fmt.Sprintf(ProjectResources, tt.arg.ID)
			projectSubResource := fmt.Sprintf(ProjectSubResources, tt.arg.ID)
			projectNameResource := fmt.Sprintf(ProjectResources, tt.arg.Name)

			authEnforcer := &enforcerMock.Enforcer{}
			if tt.authEnabled {
				authEnforcer.On(
					"UpsertRole",
					fmt.Sprintf("%s-administrators", tt.arg.Name),
					[]string(tt.arg.Administrators),
				).Return(&types.Role{
					ID:      "admin-role",
					Members: tt.arg.Administrators,
				}, nil)

				authEnforcer.On(
					"UpsertRole",
					fmt.Sprintf("%s-readers", tt.arg.Name),
					[]string(tt.arg.Readers),
				).Return(&types.Role{
					ID:      "reader-role",
					Members: tt.arg.Readers,
				}, nil)
				authEnforcer.On(
					"UpsertPolicy",
					fmt.Sprintf("%s-administrators-policy", tt.arg.Name),
					[]string{"admin-role"},
					[]string{},
					[]string{projectResource, projectSubResource, projectNameResource},
					[]string{enforcer.ActionAll},
				).Return(&types.Policy{
					ID:        "admin-policy",
					Subjects:  []string{"admin-role"},
					Resources: []string{projectResource, projectSubResource, projectNameResource},
					Actions:   []string{enforcer.ActionAll},
				}, nil)
				authEnforcer.On(
					"UpsertPolicy",
					fmt.Sprintf("%s-readers-policy", tt.arg.Name),
					[]string{"reader-role"},
					[]string{},
					[]string{projectResource, projectSubResource, projectNameResource},
					[]string{enforcer.ActionRead},
				).Return(&types.Policy{
					ID:        "reader-policy",
					Subjects:  []string{"readers-role"},
					Resources: []string{projectResource, projectSubResource, projectNameResource},
					Actions:   []string{enforcer.ActionRead},
				}, nil)
			}

			projectsService, err := NewProjectsService(MLFlowTrackingURL, storage, authEnforcer, tt.authEnabled)
			assert.NoError(t, err)

			res, err := projectsService.CreateProject(tt.arg)
			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrorMsg, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expResult, res)

			storage.AssertExpectations(t)
			authEnforcer.AssertExpectations(t)
		})
	}
}

func TestProjectsService_UpdateProject(t *testing.T) {
	tests := []struct {
		name        string
		arg         *models.Project
		authEnabled bool
		expResult   *models.Project
	}{
		{
			"success: auth enabled",
			&models.Project{
				Name:           "my-project",
				Administrators: []string{"user@email.com"},
				Readers:        nil,
			},
			true,
			&models.Project{
				Name:              "my-project",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Administrators:    []string{"user@email.com"},
				Readers:           nil,
			},
		},
		{
			"success: auth disabled",
			&models.Project{
				Name:           "my-project",
				Administrators: []string{"user@email.com"},
				Readers:        nil,
			},
			false,
			&models.Project{
				Name:              "my-project",
				MLFlowTrackingURL: MLFlowTrackingURL,
				Administrators:    []string{"user@email.com"},
				Readers:           nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mocks.ProjectRepository{}
			storage.On("Save", tt.expResult).Return(tt.expResult, nil)

			authEnforcer := &enforcerMock.Enforcer{}
			if tt.authEnabled {

				projectResource := fmt.Sprintf(ProjectResources, tt.arg.ID)
				projectSubResource := fmt.Sprintf(ProjectSubResources, tt.arg.ID)
				projectNameResource := fmt.Sprintf(ProjectResources, tt.arg.Name)

				authEnforcer.On(
					"UpsertRole",
					fmt.Sprintf("%s-administrators", tt.arg.Name),
					[]string(tt.arg.Administrators),
				).Return(&types.Role{
					ID:      "admin-role",
					Members: tt.arg.Administrators,
				}, nil)
				authEnforcer.On(
					"UpsertRole",
					fmt.Sprintf("%s-readers", tt.arg.Name),
					[]string(tt.arg.Readers),
				).Return(&types.Role{
					ID:      "reader-role",
					Members: tt.arg.Readers,
				}, nil)
				authEnforcer.On(
					"UpsertPolicy",
					fmt.Sprintf("%s-administrators-policy", tt.arg.Name),
					[]string{"admin-role"},
					[]string{},
					[]string{projectResource, projectSubResource, projectNameResource},
					[]string{enforcer.ActionAll},
				).Return(&types.Policy{
					ID:        "admin-policy",
					Subjects:  []string{"admin-role"},
					Resources: []string{projectResource, projectSubResource, projectNameResource},
					Actions:   []string{enforcer.ActionAll},
				}, nil)
				authEnforcer.On(
					"UpsertPolicy",
					fmt.Sprintf("%s-readers-policy", tt.arg.Name),
					[]string{"reader-role"},
					[]string{},
					[]string{projectResource, projectSubResource, projectNameResource},
					[]string{enforcer.ActionRead},
				).Return(&types.Policy{
					ID:        "reader-policy",
					Subjects:  []string{"readers-role"},
					Resources: []string{projectResource, projectSubResource, projectNameResource},
					Actions:   []string{enforcer.ActionRead},
				}, nil)
			}

			projectsService, err := NewProjectsService(MLFlowTrackingURL, storage, authEnforcer, tt.authEnabled)
			assert.NoError(t, err)

			res, err := projectsService.UpdateProject(tt.arg)
			assert.NoError(t, err)
			assert.Equal(t, tt.expResult, res)

			storage.AssertExpectations(t)
			authEnforcer.AssertExpectations(t)
		})
	}
}

func TestProjectsService_ListProjects(t *testing.T) {
	projectFilter := "my-project"

	exp := []*models.Project{
		{
			Name:              "my-project",
			MLFlowTrackingURL: MLFlowTrackingURL,
			Administrators:    []string{"user@email.com"},
			Readers:           nil,
		},
	}
	storage := &mocks.ProjectRepository{}
	storage.On("ListProjects", projectFilter).Return(exp, nil)

	authEnforcer := &enforcerMock.Enforcer{}
	projectsService, err := NewProjectsService(MLFlowTrackingURL, storage, authEnforcer, false)
	assert.NoError(t, err)

	res, err := projectsService.ListProjects(projectFilter)
	assert.NoError(t, err)
	assert.Equal(t, exp, res)

	storage.AssertExpectations(t)
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
	projectsService, err := NewProjectsService(MLFlowTrackingURL, storage, authEnforcer, false)
	assert.NoError(t, err)

	res, err := projectsService.FindByID(id)
	assert.NoError(t, err)
	assert.Equal(t, exp, res)

	storage.AssertExpectations(t)
}
