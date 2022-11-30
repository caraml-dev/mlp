package middleware

import (
	"net/http"
	"strings"
)

func ProjectCreationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userAgent := strings.ToLower(r.Header.Get("User-Agent"))
		if strings.Contains(userAgent, "swagger") {
			jsonError(
				w,
				"Project creation from SDK is disabled. Use the MLP console to create a project.",
				http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
