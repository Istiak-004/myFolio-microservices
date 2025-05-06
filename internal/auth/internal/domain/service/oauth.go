package service

import (
	"context"
	"fmt"

	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/ports"
)

type OAuthService struct {
	userRepo     ports.UserRepository
	tokenService ports.TokenService
	providers    map[string]ports.OAuthProvider
}

func NewOAuthService(
	userRepo ports.UserRepository,
	tokenService ports.TokenService,
	providers map[string]ports.OAuthProvider,
) *OAuthService {
	return &OAuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
		providers:    providers,
	}
}

func (s *OAuthService) Authenticate(ctx context.Context, provider, code string) (*models.TokenPair, error) {
	oauthProvider, ok := s.providers[provider]
	if !ok {
		return nil, fmt.Errorf("oauth provider %s not supported", provider)
	}

	oauthUser, err := oauthProvider.GetUserInfo(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := s.userRepo.FindOrCreateOAuthUser(ctx, provider, oauthUser.ID, oauthUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	return s.tokenService.GenerateTokens(ctx, user.ID)
}
