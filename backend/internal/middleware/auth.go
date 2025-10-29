package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/auth"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/storage"
)

// Context keys for storing user information
type contextKey string

const (
	UserContextKey contextKey = "user"
	UserIDContextKey contextKey = "user_id"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtManager *auth.JWTManager, store storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if auth is enabled
			authEnabled := os.Getenv("AUTH_ENABLED")
			if authEnabled == "false" {
				// Auth disabled - pass through without checking
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := jwtManager.ValidateAccessToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Get user from database
			user, err := store.GetUser(claims.UserID)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			ctx = context.WithValue(ctx, UserIDContextKey, user.ID)

			// Continue with user in context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the authenticated user from the request context
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	return user, ok
}

// GetUserIDFromContext retrieves the authenticated user ID from the request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}

// RequireParticipant is a middleware that checks if the user is a participant of a project
func RequireParticipant(store storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if auth is enabled
			authEnabled := os.Getenv("AUTH_ENABLED")
			if authEnabled == "false" {
				// Auth disabled - pass through
				next.ServeHTTP(w, r)
				return
			}

			// Get user ID from context
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get project ID from URL parameter (assumes chi router)
			projectID := r.Context().Value("projectID")
			if projectID == nil {
				// Try to get from chi URLParam
				projectID = r.PathValue("id")
			}

			if projectID == "" || projectID == nil {
				http.Error(w, "Project ID required", http.StatusBadRequest)
				return
			}

			// Check if user is a participant
			_, err := store.GetProjectParticipant(projectID.(string), userID)
			if err != nil {
				http.Error(w, "Access denied: not a project participant", http.StatusForbidden)
				return
			}

			// User is a participant - continue
			next.ServeHTTP(w, r)
		})
	}
}
