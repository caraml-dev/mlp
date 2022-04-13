package enforcer

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gojek/mlp/api/pkg/authz/enforcer/types"
)

// These tests are run with real keto instance running on localhost:4466
// Execute test using `make test` to spin up the necessary instances
const (
	ProductName = "test"
)

var KetoURL = getEnvOrDefault("KETO_URL", "http://localhost:4466")

var BootstrapRoles = []types.Role{
	{
		"bootrap-role-1",
		[]string{"user-1@example.com", "user-2@example.com"},
	},
	{
		"bootrap-role-2",
		[]string{"user-3@example.com", "user-4@example.com"},
	},
}

// testPolicy is a structure that holds the input data for creating a policy
type testPolicy struct {
	ID        string
	Roles     []string
	Users     []string
	Resources []string
	Actions   []string
}

var BootstrapPolicy = []testPolicy{
	{
		"bootstrap-policy-1",
		[]string{BootstrapRoles[0].ID},
		[]string{},
		[]string{"pages:1"},
		[]string{ActionRead},
	},
	{
		"bootstrap-policy-2",
		[]string{BootstrapRoles[1].ID},
		[]string{},
		[]string{"pages:**"},
		[]string{ActionAll},
	},
	{
		"bootstrap-policy-3",
		[]string{},
		[]string{"users:**"},
		[]string{"pages:10"},
		[]string{ActionRead},
	},
}

func TestEnforcer_Enforce(t *testing.T) {
	enforcer, err := NewEnforcerBuilder().URL(KetoURL).Product(ProductName).Build()
	assert.NoError(t, err)

	err = initializeBootstrapRoles(enforcer, BootstrapRoles)
	assert.NoError(t, err)
	err = initializeBootstrapPolicies(enforcer, BootstrapPolicy)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		user     string
		resource string
		action   string
		result   bool
	}{
		{
			"allow: user-1 request read pages:1",
			"user-1@example.com",
			"pages:1",
			ActionRead,
			true,
		},
		{
			"allow: user-3 request update pages:10",
			"user-3@example.com",
			"pages:10",
			ActionUpdate,
			true,
		},
		{
			"allow: user-10 request read pages:10",
			"user-10@example.com",
			"pages:10",
			ActionRead,
			true,
		},
		{
			"reject: user-10 request update pages:10",
			"user-10@example.com",
			"pages:10",
			ActionUpdate,
			false,
		},
		{
			"reject: user-2 request read pages:2",
			"user-2@example.com",
			"pages:2",
			ActionRead,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := enforcer.Enforce(tt.user, tt.resource, tt.action)
			assert.NoError(t, err)
			assert.Equal(t, tt.result, *res, "invalid enforce result")
		})
	}
}

func TestEnforcer_FilterAuthorizedResource(t *testing.T) {
	enforcer, err := NewEnforcerBuilder().URL(KetoURL).Product(ProductName).Build()
	assert.NoError(t, err)

	err = initializeBootstrapRoles(enforcer, BootstrapRoles)
	assert.NoError(t, err)
	err = initializeBootstrapPolicies(enforcer, BootstrapPolicy)
	assert.NoError(t, err)

	tests := []struct {
		name              string
		user              string
		resources         []string
		action            string
		expectedResources []string
	}{
		{
			"user-1 only able to read pages:1",
			"user-1@example.com",
			[]string{
				"pages:1",
				"pages:2",
				"pages:3",
			},
			ActionRead,
			[]string{
				"pages:1",
			},
		},
		{
			"user-3 able to update all pages",
			"user-3@example.com",
			[]string{
				"pages:1",
				"pages:2",
				"pages:3",
			},
			ActionUpdate,
			[]string{
				"pages:1",
				"pages:2",
				"pages:3",
			},
		},
		{
			"user-10 able to read pages:10",
			"user-10@example.com",
			[]string{
				"pages:1",
				"pages:10",
			},
			ActionRead,
			[]string{
				"pages:10",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := enforcer.FilterAuthorizedResource(tt.user, tt.resources, tt.action)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResources, res)
		})
	}
}

func TestEnforcer_GetRole(t *testing.T) {
	enforcer, err := NewEnforcerBuilder().URL(KetoURL).Product(ProductName).Timeout(5 * time.Second).Build()
	assert.NoError(t, err)

	err = initializeBootstrapRoles(enforcer, BootstrapRoles)
	assert.NoError(t, err)
	err = initializeBootstrapPolicies(enforcer, BootstrapPolicy)
	assert.NoError(t, err)

	role, err := enforcer.GetRole(BootstrapRoles[0].ID)
	assert.NoError(t, err)

	assert.Equal(t, "roles:test:bootrap-role-1", role.ID)
	assert.Equal(t, []string{"users:user-1@example.com", "users:user-2@example.com"}, role.Members)

	_, err = enforcer.GetRole("unknown-role")
	assert.Error(t, err)
}

func TestEnforcer_GetPolicy(t *testing.T) {
	enforcer, err := NewEnforcerBuilder().URL(KetoURL).Product(ProductName).Timeout(5 * time.Second).Build()
	assert.NoError(t, err)

	err = initializeBootstrapRoles(enforcer, BootstrapRoles)
	assert.NoError(t, err)
	err = initializeBootstrapPolicies(enforcer, BootstrapPolicy)
	assert.NoError(t, err)

	policy, err := enforcer.GetPolicy(BootstrapPolicy[0].ID)
	assert.NoError(t, err)

	assert.Equal(t, "policies:test:bootstrap-policy-1", policy.ID)
	assert.Equal(t, []string{"roles:test:bootrap-role-1"}, policy.Subjects)
	assert.Equal(t, []string{"resources:test:pages:1"}, policy.Resources)
	assert.Equal(t, []string{ActionRead}, policy.Actions)

	_, err = enforcer.GetPolicy("unknown-policy")
	assert.Error(t, err)
}

func initializeBootstrapPolicies(e Enforcer, policies []testPolicy) error {
	for _, policy := range policies {
		_, err := e.UpsertPolicy(policy.ID, policy.Roles, policy.Users, policy.Resources, policy.Actions)
		if err != nil {
			return err
		}
	}
	return nil
}

func initializeBootstrapRoles(e Enforcer, roles []types.Role) error {
	for _, role := range roles {
		_, err := e.UpsertRole(role.ID, role.Members)
		if err != nil {
			return err
		}
	}
	return nil
}

func getEnvOrDefault(env string, defaultValue string) string {
	val, ok := os.LookupEnv(env)
	if !ok {
		return defaultValue
	}
	return val
}
