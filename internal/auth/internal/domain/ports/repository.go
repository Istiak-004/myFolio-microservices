package ports

import (
	"context"

	"time"

	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/valueobjects"
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

type TokenRepository interface {
	GenerateRefreshToken(ctx context.Context, token valueobjects.Token) (string, error)
	StoreRefreshToken(ctx context.Context, userID string, token valueobjects.Token, expiresAt time.Time) error
	VerifyRefreshToken(ctx context.Context, token valueobjects.Token) (string, error)
	RevokeRefreshToken(ctx context.Context, token valueobjects.Token) error
	RotateRefreshToken(ctx context.Context, oldToken valueobjects.Token) (valueobjects.Token, error)
	RevokeAllForUser(ctx context.Context, userID string) error
	GetRefreshToken(ctx context.Context, token valueobjects.Token) (string, error)
}
