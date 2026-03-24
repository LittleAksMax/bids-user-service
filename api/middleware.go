package api

import (
	"context"
	"net/http"

	"github.com/LittleAksMax/bids-util/requests"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// RegisterMiddleware attaches common middleware to the router.
func RegisterMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
}

// InjectUUIDSubjectFromHeader parses a trusted UUID header and stores it under the standard subject context key.
func InjectUUIDSubjectFromHeader(headerKey, contextKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawUserID := r.Header.Get(headerKey)
			if rawUserID == "" {
				requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
					Success: false,
					Error:   "missing user ID header",
				})
				return
			}

			userID, err := uuid.Parse(rawUserID)
			if err != nil {
				requests.WriteJSON(w, http.StatusBadRequest, requests.APIResponse{
					Success: false,
					Error:   "invalid user ID header",
				})
				return
			}

			ctx := context.WithValue(r.Context(), contextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
