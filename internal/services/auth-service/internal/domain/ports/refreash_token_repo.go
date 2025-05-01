package ports

import (
	"context"
	"time"
)

type RefreshTokenRepository interface {
	Store(ctx context.Context, userID, token string, expiresAt time.Time) error
	Get(ctx context.Context, token string) (string, time.Time, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}
