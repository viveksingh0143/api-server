package middlewares

import (
	"net/http"
)

// ContentTypeMiddleware sets the Content-Type header to application/json
func ContentTypeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
