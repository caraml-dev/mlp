package enforcer

import (
	"context"
	"fmt"
	"sort"
	"testing"

	ory "github.com/ory/keto-client-go"
	"github.com/stretchr/testify/require"
)

const (
	ketoRemoteRead  = "http://localhost:4466"
	ketoRemoteWrite = "http://localhost:4467"
)

// These tests are run with real keto instance running on localhost, with port 4466 and 4467 for read and write
// respectively. Execute test using `make test` to spin up the necessary instances
func TestNewEnforcer(t *testing.T) {
	tests := map[string]struct {
		ketoRemoteRead  string
		ketoRemoteWrite string
		cacheConfig     *CacheConfig

		expectedError string
	}{
		"success | no cache": {
			ketoRemoteRead:  "http://localhost:4466",
			ketoRemoteWrite: "http://localhost:4467",
		},
		"success | with cache": {
			ketoRemoteRead:  "http://localhost:4466",
			ketoRemoteWrite: "http://localhost:4467",
			cacheConfig: &CacheConfig{
				KeyExpirySeconds:            30,
				CacheCleanUpIntervalSeconds: 60,
			},
		},
		"failure | large cache expiry": {
			ketoRemoteRead:  "http://localhost:4466",
			ketoRemoteWrite: "http://localhost:4467",
			cacheConfig: &CacheConfig{
				KeyExpirySeconds:            3000,
				CacheCleanUpIntervalSeconds: 60,
			},
			expectedError: "Configured KeyExpirySeconds is larger than the max permitted value of 600",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			builder := NewEnforcerBuilder().
				KetoEndpoints(tt.ketoRemoteRead, tt.ketoRemoteWrite)
			if tt.cacheConfig != nil {
				builder.WithCaching(tt.cacheConfig.KeyExpirySeconds, tt.cacheConfig.CacheCleanUpIntervalSeconds)
			}
			_, err := builder.Build()

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEnforcer_HasPermission(t *testing.T) {
	ketoEnforcer, err := NewEnforcerBuilder().Build()
	require.NoError(t, err)
	readClient := newKetoClient(ketoRemoteRead)
	writeClient := newKetoClient(ketoRemoteWrite)
	clearRelations(readClient, writeClient)

	newRoleAndPermissionsRequest := NewAuthorizationUpdateRequest()
	newRoleAndPermissionsRequest.SetRolePermissions("page.1.admin", []string{"page.1.get", "page.1.put"})
	newRoleAndPermissionsRequest.SetRoleMembers("page.1.admin", []string{"user-1@example.com"})
	err = ketoEnforcer.UpdateAuthorization(context.Background(), newRoleAndPermissionsRequest)
	require.NoError(t, err)

	tests := []struct {
		name       string
		permission string
		user       string
		result     bool
	}{
		{
			"allow: user-1 request read page.1",
			"page.1.get",
			"user-1@example.com",
			true,
		},
		{
			"reject: user-3 request update page.1",
			"page.1.put",
			"user-3@example.com",
			false,
		},
		{
			"reject: user-1 request read unknown page number",
			"page.99.put",
			"user-1@example.com",
			false,
		},
		{
			"reject: unknown user request read page 1",
			"page.1.put",
			"anonymous@example.com",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ketoEnforcer.IsUserGrantedPermission(context.Background(), tt.user, tt.permission)
			require.NoError(t, err)
			require.Equal(t, tt.result, res, "invalid enforce result")
		})
	}
	updateRoleAndPermissionsRequest := NewAuthorizationUpdateRequest()
	updateRoleAndPermissionsRequest.SetRolePermissions("page.1.admin", []string{"page.1.get", "page.1.delete"})
	updateRoleAndPermissionsRequest.SetRoleMembers("page.1.admin", []string{"admin-1@example.com"})
	err = ketoEnforcer.UpdateAuthorization(context.Background(), updateRoleAndPermissionsRequest)
	testsAfterUpdate := []struct {
		name       string
		permission string
		user       string
		result     bool
	}{
		{
			"reject after update: user-1 request read page.1",
			"page.1.get",
			"user-1@example.com",
			false,
		},
		{
			"allow after update: admin-1 request delete page.1",
			"page.1.delete",
			"admin-1@example.com",
			true,
		},
		{
			"reject after update: admin-1 request update page.1",
			"page.1.put",
			"admin-1@example.com",
			false,
		},
	}
	for _, tt := range testsAfterUpdate {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ketoEnforcer.IsUserGrantedPermission(context.Background(), tt.user, tt.permission)
			require.NoError(t, err)
			require.Equal(t, tt.result, res, "invalid enforce result")
		})
	}

	require.NoError(t, err)
}

func TestEnforcer_GetUserRoles(t *testing.T) {
	ketoEnforcer, err := NewEnforcerBuilder().Build()
	require.NoError(t, err)
	readClient := newKetoClient(ketoRemoteRead)
	writeClient := newKetoClient(ketoRemoteWrite)
	clearRelations(readClient, writeClient)
	updateRequest := NewAuthorizationUpdateRequest()
	for i := 1; i < 4; i++ {
		updateRequest.SetRoleMembers(fmt.Sprintf("pages.%d.reader", i), []string{"user-1@example.com"})
	}
	err = ketoEnforcer.UpdateAuthorization(context.Background(), updateRequest)
	require.NoError(t, err)

	tests := []struct {
		name          string
		user          string
		expectedRoles []string
	}{
		{
			"user-1 is reader for page 1, 2 and 3",
			"user-1@example.com",
			[]string{
				"pages.1.reader",
				"pages.2.reader",
				"pages.3.reader",
			},
		},
		{
			"unknown user has no roles",
			"anonymous@example.com",
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ketoEnforcer.GetUserRoles(context.Background(), tt.user)
			require.NoError(t, err)
			sort.Strings(tt.expectedRoles)
			sort.Strings(res)
			require.Equal(t, tt.expectedRoles, res)
		})
	}
}

func TestEnforcer_GetRolePermissions(t *testing.T) {
	ketoEnforcer, err := NewEnforcerBuilder().Build()
	require.NoError(t, err)
	readClient := newKetoClient(ketoRemoteRead)
	writeClient := newKetoClient(ketoRemoteWrite)
	clearRelations(readClient, writeClient)
	updateRequest := NewAuthorizationUpdateRequest()
	updateRequest.SetRolePermissions("pages.1.reader", []string{"pages.1.get", "pages.1.post"})
	err = ketoEnforcer.UpdateAuthorization(context.Background(), updateRequest)
	require.NoError(t, err)
	tests := []struct {
		name                string
		role                string
		expectedPermissions []string
	}{
		{
			"pages.1.reader role can read and update page 1",
			"pages.1.reader",
			[]string{
				"pages.1.get",
				"pages.1.post",
			},
		},
		{
			"unknown user has no roles",
			"anonymous@example.com",
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ketoEnforcer.GetRolePermissions(context.Background(), tt.role)
			require.NoError(t, err)
			sort.Strings(tt.expectedPermissions)
			sort.Strings(res)
			require.Equal(t, tt.expectedPermissions, res)
		})
	}
}

func TestEnforcer_GetRoleMembers(t *testing.T) {
	ketoEnforcer, err := NewEnforcerBuilder().Build()
	require.NoError(t, err)
	readClient := newKetoClient(ketoRemoteRead)
	writeClient := newKetoClient(ketoRemoteWrite)
	clearRelations(readClient, writeClient)
	updateRequest := NewAuthorizationUpdateRequest()
	updateRequest.SetRolePermissions("pages.1.reader", []string{"pages.1.get", "pages.1.post"})
	err = ketoEnforcer.UpdateAuthorization(context.Background(), updateRequest)
	require.NoError(t, err)
	tests := []struct {
		name                string
		role                string
		expectedPermissions []string
	}{
		{
			"pages.1.reader role can read and update page 1",
			"pages.1.reader",
			[]string{
				"pages.1.get",
				"pages.1.post",
			},
		},
		{
			"unknown user has no roles",
			"anonymous@example.com",
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ketoEnforcer.GetRolePermissions(context.Background(), tt.role)
			require.NoError(t, err)
			sort.Strings(tt.expectedPermissions)
			sort.Strings(res)
			require.Equal(t, tt.expectedPermissions, res)
		})
	}
}

func newKetoClient(endpoint string) *ory.APIClient {
	cfg := ory.NewConfiguration()
	cfg.Servers = ory.ServerConfigurations{
		{
			URL: endpoint,
		},
	}
	return ory.NewAPIClient(cfg)
}
func clearRelations(readClient *ory.APIClient, writeClient *ory.APIClient) {
	ctx := context.Background()
	resp, _, err := readClient.RelationshipApi.ListRelationshipNamespaces(ctx).Execute()
	if err != nil {
		panic(err)
	}
	for _, namespace := range resp.GetNamespaces() {
		_, err := writeClient.RelationshipApi.DeleteRelationships(context.Background()).Namespace(namespace.GetName()).
			Execute()
		if err != nil {
			panic(err)
		}
	}
}
