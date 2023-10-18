package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
)

func TestStartKetoBootstrap(t *testing.T) {
	bootstrapCfg := &BootstrapConfig{
		KetoRemoteRead:  "http://localhost:4466",
		KetoRemoteWrite: "http://localhost:4467",
		ProjectReaders:  []string{"user1", "user2"},
		MLPAdmins:       []string{"admin1", "admin2"},
	}

	// Create a mock enforcer
	mockEnforcer := &enforcer.Enforcer{
		UpdateAuthorizationFunc: func(ctx context.Context, req *enforcer.AuthorizationUpdateRequest) error {
			// Check that the role members and permissions are set correctly
			assert.Equal(t, bootstrapCfg.ProjectReaders, req.RoleMembers(enforcer.MLPProjectsReaderRole))
			assert.Equal(t, bootstrapCfg.MLPAdmins, req.RoleMembers(enforcer.MLPAdminRole))
			assert.Equal(t, []string{"mlp.projects.post"}, req.RolePermissions(enforcer.MLPAdminRole))
			return nil
		},
	}

	// Replace the NewEnforcerBuilder function with a mock that returns the mock enforcer
	enforcer.NewEnforcerBuilder = func() *enforcer.EnforcerBuilder {
		return &enforcer.EnforcerBuilder{
			Enforcer: mockEnforcer,
		}
	}

	err := startKetoBootstrap(bootstrapCfg)
	require.NoError(t, err)
}
