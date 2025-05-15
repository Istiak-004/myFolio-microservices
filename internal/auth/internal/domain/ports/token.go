package ports

import (
	"context"

	"github.com/istiak-004/myFolio-microservices/auth/internal/infrastructure/token"
)

type JWTService interface {
	GenerateAccessToken(userID, role string) (signed string, jti string, err error)
	GenerateRefreshToken(userID string) (signed string, jti string, err error)
	VerifyAccessToken(ctx context.Context, tokenStr string) (*token.CustomClaims, error)
}
