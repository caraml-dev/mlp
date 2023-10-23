package cmd

import (
	"context"
	"testing"

	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
	"github.com/stretchr/testify/require"
)

func TestStartKetoBootsrap(t *testing.T) {
	bootstrapCfg := &BootstrapConfig{
		KetoRemoteRead:  "http://localhost:4466",
		KetoRemoteWrite: "http://localhost:4467",
		ProjectReaders:  []string{},
		MLPAdmins:       []string{"admin@email.com"}}

	tests := []struct {
		name            string
		permission      string
		user            string
		bootstrapEnable bool
		result          bool
	}{
		{
			"disable: user-1 request create project",
			"mlp.projects.post",
			"user-1@example.com",
			true,
			false, // user-1 can't create project
		},
		{
			"allow: admin request create project",
			"mlp.projects.post",
			"admin@email.com",
			true,
			true, // admin can create project
		},
		{
			"allow: user-1 request create project",
			"mlp.projects.post",
			"user-1@example.com",
			false,
			true, // user-1 can't create project
		},
		{
			"allow: admin request create project",
			"mlp.projects.post",
			"admin@email.com",
			false,
			true, // admin can create project
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// if bootstrap enabled run bootstrap
			if tt.bootstrapEnable {
				err := startKetoBootstrap(bootstrapCfg)
				require.NoError(t, err)
			} else {
				// else updateRequest := enforcer.NewAuthorizationUpdateRequest()
				ketoEnforcer, err := enforcer.NewEnforcerBuilder().Build()
				require.NoError(t, err)

				newRoleAndPermissionsRequest := enforcer.NewAuthorizationUpdateRequest()
				newRoleAndPermissionsRequest.AddRolePermissions("page.1.admin", []string{"mlp.projects.post"})
				newRoleAndPermissionsRequest.SetRoleMembers("page.1.admin", []string{tt.user})
				err = ketoEnforcer.UpdateAuthorization(context.Background(), newRoleAndPermissionsRequest)
				require.NoError(t, err)

			}

			authEnforcer, err := enforcer.NewEnforcerBuilder().
				KetoEndpoints(bootstrapCfg.KetoRemoteRead, bootstrapCfg.KetoRemoteWrite).
				Build()
			require.NoError(t, err)

			res, err := authEnforcer.IsUserGrantedPermission(context.Background(), tt.user, tt.permission)
			require.NoError(t, err)
			require.Equal(t, tt.result, res)
		})
	}
}
