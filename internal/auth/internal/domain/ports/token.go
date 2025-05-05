package ports

import (
	"context"
	"time"
)

type TokenService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken() (string, error)
	VerifyAccessToken(tokenString string) (string, error)
	StoreRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	VerifyRefreshToken(ctx context.Context, token string) (string, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RotateRefreshToken(ctx context.Context, oldToken string) (newToken string, err error)
}
