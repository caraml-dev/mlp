package cmd

import (
	"testing"

	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
	enforcerMock "github.com/caraml-dev/mlp/api/pkg/authz/enforcer/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStartKetoBootsrap(t *testing.T) {
	tests := []struct {
		name                               string
		projectReaders                     []string
		mlpAdmins                          []string
		expectedUpdateAuthorizationRequest enforcer.AuthorizationUpdateRequest
	}{
		{
			"admin role must have project post permission even there are no project readers",
			[]string{},
			[]string{"admin1"},
			enforcer.AuthorizationUpdateRequest{
				RolePermissions: map[string][]string{
					"mlp.administrator": {"mlp.projects.post"},
				},
				RoleMembers: map[string][]string{
					"mlp.projects.reader": {},
					"mlp.administrator":   {"admin1"},
				},
			},
		},
		{
			"admin role should have project post permission, even there are no mlp admins or project readers",
			[]string{},
			[]string{},
			enforcer.AuthorizationUpdateRequest{
				RolePermissions: map[string][]string{
					"mlp.administrator": {"mlp.projects.post"},
				},
				RoleMembers: map[string][]string{
					"mlp.projects.reader": {},
					"mlp.administrator":   {},
				},
			},
		},
		{
			"only admin role should have project post permission, even there are no mlp admins and even there are project readers",
			[]string{"readers1", "readers2"},
			[]string{},
			enforcer.AuthorizationUpdateRequest{
				RolePermissions: map[string][]string{
					"mlp.administrator": {"mlp.projects.post"},
				},
				RoleMembers: map[string][]string{
					"mlp.projects.reader": {"readers1", "readers2"},
					"mlp.administrator":   {},
				},
			},
		},
		{
			"only admin role should have project post permission, even there are project readers",
			[]string{"readers1", "readers2"},
			[]string{"admin1"},
			enforcer.AuthorizationUpdateRequest{
				RolePermissions: map[string][]string{
					"mlp.administrator": {"mlp.projects.post"},
				},
				RoleMembers: map[string][]string{
					"mlp.projects.reader": {"readers1", "readers2"},
					"mlp.administrator":   {"admin1"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authEnforcer := &enforcerMock.Enforcer{}

			authEnforcer.On("UpdateAuthorization", mock.Anything, tt.expectedUpdateAuthorizationRequest).Return(nil)
			err := startKetoBootstrap(authEnforcer, tt.projectReaders, tt.mlpAdmins)
			authEnforcer.AssertExpectations(t)
			require.NoError(t, err)
		})
	}
}
