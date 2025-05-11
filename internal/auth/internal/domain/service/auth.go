package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/istiak-004/myFolio-microservices/auth/config"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/ports"
	"github.com/istiak-004/myFolio-microservices/auth/pkg/security"
	"github.com/istiak-004/myFolio-microservices/auth/pkg/token"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotVerified    = errors.New("user not verified")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenInvalid       = errors.New("invalid token")
)

type AuthServiceImpl struct {
	userRepo         ports.UserRepository
	verificationRepo ports.VerificationRepository
	tokenRepo        ports.TokenRepository
	accessTokenGen   *token.JWTGenerator
	refreshTokenGen  *token.JWTGenerator
	mailer           ports.Mailer
	appConfig        *config.Config
}

func NewAuthService(
	userRepo ports.UserRepository,
	verificationRepo ports.VerificationRepository,
	tokenRepo ports.TokenRepository,
	accessTokenGen *token.JWTGenerator,
	refreshTokenGen *token.JWTGenerator,
	mailer ports.Mailer,
	appConfig *config.Config,
) ports.AuthService {
	return &AuthServiceImpl{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		tokenRepo:        tokenRepo,
		accessTokenGen:   accessTokenGen,
		refreshTokenGen:  refreshTokenGen,
		mailer:           mailer,
		appConfig:        appConfig,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, email, password, firstName, lastName string) (*models.User, error) {
	// Check if email exists
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: hashedPassword,
		IsVerified:   false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	randKey, err := security.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	// Generate verification token
	verificationToken := models.VerificationToken{
		Token:     randKey,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.verificationRepo.Create(ctx, &verificationToken); err != nil {
		return nil, err
	}

	// Send verification email
	verificationURL := fmt.Sprintf("%s%s?token=%s",
		s.appConfig.App.BaseURL,
		s.appConfig.App.VerificationPath,
		verificationToken.Token)

	if err := s.mailer.SendVerificationEmail(user.Email, firstName+" "+lastName, verificationURL); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*models.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, ErrUserNotVerified
	}

	if err := security.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokens(ctx, user.ID)
}

func (s *AuthServiceImpl) VerifyEmail(ctx context.Context, token string) error {
	verification, err := s.verificationRepo.Get(ctx, token)
	if err != nil {
		return err
	}
	if verification == nil {
		return ErrTokenInvalid
	}

	if time.Now().After(verification.ExpiresAt) {
		return ErrTokenExpired
	}

	user, err := s.userRepo.FindByID(ctx, verification.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.IsVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.verificationRepo.Delete(ctx, token)
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	userID, err := s.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Revoke old token
	if err := s.tokenRepo.RevokeRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, userID)
}

func (s *AuthServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenRepo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *AuthServiceImpl) generateTokens(ctx context.Context, userID string) (*models.TokenPair, error) {
	accessToken, err := s.accessTokenGen.Generate(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshTokenGen.Generate(userID)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.refreshTokenGen.GetExpiry())
	if err := s.tokenRepo.StoreRefreshToken(ctx, userID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.accessTokenGen.GetExpiry().Seconds()),
	}, nil
}
