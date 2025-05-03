package service

import (
	"context"
	"errors"
	"time"

	"myFolio-microservices/internal/services/myFolio-microservices/internal/domain/models"
	"myFolio-microservices/internal/services/myFolio-microservices/internal/domain/ports"
	"myFolio-microservices/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenInvalid       = errors.New("invalid token")
)

type AuthServiceImpl struct {
	userRepo      ports.UserRepository
	sessionRepo   ports.SessionRepository
	oauthProvider ports.OAuthProvider
	tokenService  ports.TokenService
	tokenExpiry   time.Duration
}

func NewAuthService(
	userRepo ports.UserRepository,
	sessionRepo ports.SessionRepository,
	oauthProvider ports.OAuthProvider,
	tokenService ports.TokenService,
	tokenExpiry time.Duration,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		oauthProvider: oauthProvider,
		tokenService:  tokenService,
		tokenExpiry:   tokenExpiry,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, email, password string) (*models.User, error) {
	if !utils.ValidateEmail(email) || !utils.ValidatePassword(password) {
		return nil, ErrInvalidCredentials
	}

	// check is user exist or not
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           utils.GenerateUUID(),
		Email:        email,
		PasswordHash: hashedPassword,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*models.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return s.tokenService.GenerateTokenPair(ctx, user.ID)

}
