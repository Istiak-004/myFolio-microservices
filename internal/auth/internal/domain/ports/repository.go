package ports

import (
	"context"

	"time"

	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	FindByGitHubID(ctx context.Context, githubID string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

// VerificationRepository defines the interface for verification token persistence
type VerificationRepository interface {
	Create(ctx context.Context, token *models.VerificationToken) error
	Get(ctx context.Context, token string) (*models.VerificationToken, error)
	Delete(ctx context.Context, token string) error
}

// type RefreshTokenRepository interface {
// 	Store(ctx context.Context, userID, token string, expiresAt time.Time) error
// 	Get(ctx context.Context, token string) (string, time.Time, error)
// 	Revoke(ctx context.Context, token string) error
// 	RevokeAllForUser(ctx context.Context, userID string) error
// }

type TokenRepository interface {
	// GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(ctx context.Context, token string) (string, error)
	// VerifyAccessToken(tokenString string) (string, error)
	StoreRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	VerifyRefreshToken(ctx context.Context, token string) (string, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RotateRefreshToken(ctx context.Context, oldToken string) (newToken string, err error)
	RevokeAllForUser(ctx context.Context, userID string) error
	GetRefreshToken(ctx context.Context, token string) (string, error)
}
