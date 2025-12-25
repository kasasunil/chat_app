package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kasasunil/chat_app/config"
	"github.com/kasasunil/chat_app/internal/pkg/errors"
	"github.com/kasasunil/chat_app/internal/pkg/logger"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware provides authentication for routes
type AuthMiddleware struct {
	config *config.Config
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

// Authenticate is the middleware function that validates authentication
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get(HeaderAuthorization)
		if authHeader == "" {
			logger.Warn(logger.TraceAuthHeaderMissing)
			respondWithError(w, errors.ErrAuthRequired)
			return
		}

		// Check if it's Basic auth
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != AuthSchemeBasic {
			logger.Warn(logger.TraceAuthInvalidFormat)
			respondWithError(w, errors.ErrInvalidAuthFormat)
			return
		}

		// Decode base64 credentials
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			logger.Warn(logger.TraceAuthInvalidBase64)
			respondWithError(w, errors.ErrInvalidBase64)
			return
		}

		// Split username and password
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			logger.Warn(logger.TraceAuthInvalidCreds)
			respondWithError(w, errors.ErrInvalidCredentialsFormat)
			return
		}

		username := credentials[0]
		password := credentials[1]

		// Iterate through all auth clients and check if credentials match
		authenticated := false
		for _, client := range m.config.GetAuthClients() {
			if client.Username == username && client.Password == password {
				authenticated = true
				break
			}
		}

		if !authenticated {
			logger.Warn(logger.TraceAuthFailed, username)
			respondWithError(w, errors.ErrInvalidCredentials)
			return
		}

		logger.Debug(logger.TraceAuthSuccess, username)
		// Add username to context (can be used as user identifier)
		ctx := context.WithValue(r.Context(), UserIDKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts user ID from request context
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// respondWithError sends an error response using the common error format
func respondWithError(w http.ResponseWriter, err *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatus)
	json.NewEncoder(w).Encode(err.ToJSONResponse())
}
