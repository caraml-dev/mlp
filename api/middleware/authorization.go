package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
)

type Authorizer struct {
	authEnforcer enforcer.Enforcer
}

func NewAuthorizer(enforcer enforcer.Enforcer) *Authorizer {
	return &Authorizer{authEnforcer: enforcer}
}

type Operation struct {
	RequestPath   string
	RequestMethod []string
}

var publicOperations = []Operation{
	{"/projects", []string{http.MethodGet, http.MethodPost}},
	{"/applications", []string{http.MethodGet}},
}

// AuthorizationMiddleware is a middleware that checks if the request is authorized.
func (a *Authorizer) AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.RequireAuthorization(r.URL.Path, r.Method) {
			next.ServeHTTP(w, r)
			return
		}
		permission := a.GetPermission(r.URL.Path, r.Method)
		user := r.Header.Get("User-Email")

		allowed, err := a.authEnforcer.IsUserGrantedPermission(r.Context(), user, permission)
		if err != nil {
			jsonError(
				w,
				fmt.Sprintf("Error while checking authorization: %s", err),
				http.StatusInternalServerError)
			return
		}
		if !allowed {
			jsonError(
				w,
				fmt.Sprintf("%s does not have the permission:%s ", user, permission),
				http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireAuthorization returns true if the request requires authorization.
func (a *Authorizer) RequireAuthorization(requestPath string, requestMethod string) bool {
	if requestMethod == http.MethodOptions {
		return false
	}

	for _, operation := range publicOperations {
		if operation.RequestPath == requestPath && slices.Contains(operation.RequestMethod, requestMethod) {
			return false
		}
	}
	return true
}

// GetPermission returns the permission required to authorized a request.
// It's assumed that permission to request is a one to one mapping.
func (a *Authorizer) GetPermission(requestPath string, requestMethod string) string {
	parts := strings.Split(strings.TrimPrefix(requestPath, "/"), "/")
	// Current paths registered in MLP are of the following format:
	// - /projects
	// - /applications
	// - /projects/{project_id}/**
	// Only project sub-resources endpoint require permission. If a user has READ/WRITE permissions
	// on /projects/{project_id}, they would also have the same permissions on all its sub-resources.
	if len(parts) > 1 {
		parts = parts[:2]
	}
	return fmt.Sprintf("mlp.%s.%s", strings.Join(parts, "."), strings.ToLower(requestMethod))
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if len(msg) > 0 {
		errJSON, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{msg})

		_, _ = w.Write(errJSON)
	}
}
