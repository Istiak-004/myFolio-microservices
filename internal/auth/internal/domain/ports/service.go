package ports

import (
	"context"

	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
)

// AuthService defines the core authentication service interface
type AuthService interface {
	Register(ctx context.Context, email, password, name string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.TokenPair, error)
	VerifyEmail(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type OAuthService interface {
	AuthURL(provider string) string
	HandleCallback(ctx context.Context, provider, code string) (*models.TokenPair, error)
}

type Mailer interface {
	SendVerificationEmail(email, name, verificationURL string) error
}
