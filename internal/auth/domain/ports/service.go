package ports

import (
	"myFolio-microservices/internal/services/myFolio-microservices/internal/domain/models"
	"context"
)

// AuthService defines the core authentication service interface
type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	VerifyToken(ctx context.Context, token string) (string, error)
	GoogleAuth(ctx context.Context, code string) (*models.TokenPair, error)
}

// OAuthProvider defines the interface for OAuth providers
type OAuthProvider interface {
	GetUserInfo(ctx context.Context, code string) (*models.OAuthUser, error)
}
