package auth_service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/ports"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/valueobjects"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotVerified    = errors.New("user not verified")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenInvalid       = errors.New("invalid token")
)

type authService struct {
	users         ports.UserRepository
	verifications ports.VerificationRepository
	tokens        ports.TokenRepository
	jwt           ports.JWTService
	mailer        ports.Mailer
}

func NewAuthService(
	users ports.UserRepository,
	verifications ports.VerificationRepository,
	tokens ports.TokenRepository,
	jwt ports.JWTService,
	mailer ports.Mailer,
) ports.AuthService {
	return &authService{users, verifications, tokens, jwt, mailer}
}

func (s *authService) Register(ctx context.Context, email valueobjects.Email, password valueobjects.Password, name string) (*models.User, error) {
	hash, err := password.Hash()
	if err != nil {
		return nil, err
	}

	existingUser, err := s.users.FindByEmail(ctx, email.String())
	if existingUser != nil || err != nil {
		return nil, ErrEmailExists
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email.String(),
		PasswordHash: hash,
		IsActive:     true,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		FirstName:    name,
		Role:         "visitor",
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	token := valueobjects.NewToken()
	verif := &models.VerificationToken{
		Token:     token.String(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.verifications.Create(ctx, verif); err != nil {
		return nil, err
	}
	_ = s.mailer.SendVerificationEmail(user.Email, user.FirstName, token.String())
	return user, nil
}

func (s *authService) Login(ctx context.Context, email valueobjects.Email, password valueobjects.Password) (*models.TokenPair, error) {
	user, err := s.users.FindByEmail(ctx, email.String())
	if err != nil {
		return nil, err
	}
	if !password.Matches(user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	access, jti, err := s.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	refreshToken := valueobjects.NewTokenWithJTI(jti)
	refresh, _, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	if err := s.tokens.StoreRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(7*24*time.Hour)); err != nil {
		return nil, err
	}
	return &models.TokenPair{AccessToken: access, RefreshToken: refresh, ExpiresIn: 3600}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, token valueobjects.Token) error {
	verif, err := s.verifications.Get(ctx, token.String())
	if err != nil {
		return err
	}
	if verif.ExpiresAt.Before(time.Now()) {
		return errors.New("verification token expired")
	}
	user, err := s.users.FindByID(ctx, verif.UserID)
	if err != nil {
		return err
	}
	user.IsVerified = true
	user.UpdatedAt = time.Now()
	if err := s.users.Update(ctx, user); err != nil {
		return err
	}
	_ = s.verifications.Delete(ctx, token.String())
	return nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken valueobjects.Token) (*models.TokenPair, error) {
	userID, err := s.tokens.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	newToken, err := s.tokens.RotateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	access, _, err := s.jwt.GenerateAccessToken(userID, "visitor")
	if err != nil {
		return nil, err
	}
	return &models.TokenPair{AccessToken: access, RefreshToken: newToken.String(), ExpiresIn: 3600}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken valueobjects.Token) error {
	return s.tokens.RevokeRefreshToken(ctx, refreshToken)
}
