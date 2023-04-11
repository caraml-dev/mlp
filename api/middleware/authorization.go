package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
)

func NewAuthorizer(enforcer enforcer.Enforcer) *Authorizer {
	return &Authorizer{authEnforcer: enforcer}
}

type Authorizer struct {
	authEnforcer enforcer.Enforcer
}

var methodMapping = map[string]string{
	http.MethodGet:     enforcer.ActionRead,
	http.MethodHead:    enforcer.ActionRead,
	http.MethodPost:    enforcer.ActionCreate,
	http.MethodPut:     enforcer.ActionUpdate,
	http.MethodPatch:   enforcer.ActionUpdate,
	http.MethodDelete:  enforcer.ActionDelete,
	http.MethodConnect: enforcer.ActionRead,
	http.MethodOptions: enforcer.ActionRead,
	http.MethodTrace:   enforcer.ActionRead,
}

func (a *Authorizer) AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resource, err := a.getResource(r)
		if err != nil {
			jsonError(
				w,
				fmt.Sprintf("Error while checking authorization: %s", err),
				http.StatusInternalServerError)
			return
		}

		action := methodToAction(r.Method)
		user := r.Header.Get("User-Email")

		allowed, err := a.authEnforcer.Enforce(user, resource, action)
		if err != nil {
			jsonError(
				w,
				fmt.Sprintf("Error while checking authorization: %s", err),
				http.StatusInternalServerError)
			return
		}
		if !*allowed {
			jsonError(
				w,
				fmt.Sprintf("%s is not authorized to execute %s on %s", user, action, resource),
				http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authorizer) getResource(r *http.Request) (string, error) {
	resource := strings.Replace(strings.TrimPrefix(r.URL.Path, "/"), "/", ":", -1)
	return resource, nil
}

func methodToAction(method string) string {
	return methodMapping[method]
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
