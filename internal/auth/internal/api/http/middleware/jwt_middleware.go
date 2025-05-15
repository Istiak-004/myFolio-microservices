package middleware

import (
	"context"

	"net/http"
	"strings"

	"github.com/istiak-004/myFolio-microservices/auth/internal/infrastructure/token"
)

type contextKey string

const (
	ContextUserIDKey contextKey = "userID"
	ContextRoleKey   contextKey = "role"
)

type JWTMiddleware struct {
	TokenVerifier token.JWTTokenManager
}

func NewJWTMiddleware(verifier token.JWTTokenManager) *JWTMiddleware {
	return &JWTMiddleware{
		TokenVerifier: verifier,
	}
}

func (m *JWTMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or malformed token", http.StatusUnauthorized)
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.TokenVerifier.VerifyAccessToken(r.Context(), rawToken)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserIDKey, claims.Subject)
		ctx = context.WithValue(ctx, ContextRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(ContextUserIDKey).(string)
	return val, ok
}

// GetUserRoleFromContext extracts role from context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(ContextRoleKey).(string)
	return val, ok
}
