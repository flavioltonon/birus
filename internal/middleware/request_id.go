package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

// SetRequestID adds a request-id header to the request
func SetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString()
		r.Header.Add("request-id", requestID)
		w.Header().Add("request-id", requestID)
		next.ServeHTTP(w, r)
	})
}
