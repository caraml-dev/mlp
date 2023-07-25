package enforcer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCache_LookUpUserPermission(t *testing.T) {
	cache := newInMemoryCache(600, 600)
	cache.StoreUserPermission("user1@email.com", "mlp.projects.1.get", true)
	cache.StoreUserPermission("user1@email.com", "mlp.projects.1.post", false)
	trueValue, falseValue := true, false
	tests := map[string]struct {
		user          string
		permission    string
		expectedVal   *bool
		expectedFound bool
	}{
		"cache hit | both user and permission match, result is true": {
			user:          "user1@email.com",
			permission:    "mlp.projects.1.get",
			expectedVal:   &trueValue,
			expectedFound: true,
		},
		"cache hit | both user and permission match, result is false": {
			user:          "user1@email.com",
			permission:    "mlp.projects.1.post",
			expectedVal:   &falseValue,
			expectedFound: true,
		},
		"cache miss | user matches but permission doesn't": {
			user:          "user1@email.com",
			permission:    "mlp.projects.1.delete",
			expectedVal:   nil,
			expectedFound: false,
		},
		"cache miss | permission matches but user doesn't": {
			user:          "user2@email.com",
			permission:    "mlp.projects.1.get",
			expectedVal:   nil,
			expectedFound: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cachedVal, found := cache.LookUpUserPermission(tt.user, tt.permission)
			assert.Equal(t, tt.expectedVal, cachedVal)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}
