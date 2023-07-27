package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"

	enforcerMock "github.com/caraml-dev/mlp/api/pkg/authz/enforcer/mocks"
)

func TestAuthorizer_RequireAuthorization(t *testing.T) {
	authorizer := NewAuthorizer(&enforcerMock.Enforcer{})
	tests := []struct {
		name     string
		path     string
		method   string
		expected bool
	}{
		{"All authenticated users can list projects", "/projects", "GET", false},
		{"All authenticated users can create new project", "/projects", "POST", false},
		{"All authenticated users can list applications", "/applications", "GET", false},
		{"Only authorized users can update project", "/projects/100", "PATCH", true},
		{"Options http request does not require authorization", "/projects/100", "OPTIONS", false},
		{"Only authorized users can access project sub resources", "/projects/100/secrets", "GET", true},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, authorizer.RequireAuthorization(tt.path, tt.method))
	}
}

func TestAuthorizer_GetPermission(t *testing.T) {
	authorizer := NewAuthorizer(&enforcerMock.Enforcer{})
	tests := []struct {
		name     string
		path     string
		method   string
		expected string
	}{
		{"project permission", "/projects/1003", "GET", "mlp.projects.1003.get"},
		{"project sub-resource permission", "/projects/1003/secrets", "GET", "mlp.projects.1003.get"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, authorizer.GetPermission(tt.path, tt.method))
	}
}
