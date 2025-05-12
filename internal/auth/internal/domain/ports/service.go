package ports

import (
	"context"

	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/valueobjects"
)

// AuthService defines the core authentication service interface
type AuthService interface {
	Register(ctx context.Context, email valueobjects.Email, password valueobjects.Password, name string) (*models.User, error)
	Login(ctx context.Context, email valueobjects.Email, password valueobjects.Password) (*models.TokenPair, error)
	VerifyEmail(ctx context.Context, token valueobjects.Token) error
	RefreshToken(ctx context.Context, refreshToken valueobjects.Token) (*models.TokenPair, error)
	Logout(ctx context.Context, refreshToken valueobjects.Token) error
}

type OAuthService interface {
	AuthURL(provider string) string
	HandleCallback(ctx context.Context, provider, code string) (*models.TokenPair, error)
}

type Mailer interface {
	SendVerificationEmail(email, name, verificationURL string) error
}
