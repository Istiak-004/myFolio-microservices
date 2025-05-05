package token

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/istiak-004/myFolio-microservices/auth/domain/ports"
)

type JWTTokenService struct {
	private          *rsa.PrivateKey
	public           *rsa.PublicKey
	issuer           string
	accessTokenExp   time.Duration
	refreshTokenExp  time.Duration
	refreshTokenRepo ports.RefreshTokenRepository
}

func NewJWTService(
	private *rsa.PrivateKey,
	public *rsa.PublicKey,
	issuer string,
	accessTokenExp time.Duration,
	refreshTokenExp time.Duration,
	refreshTokenRepo ports.RefreshTokenRepository,
) *JWTTokenService {
	return &JWTTokenService{
		private:          private,
		public:           public,
		issuer:           issuer,
		accessTokenExp:   accessTokenExp,
		refreshTokenExp:  refreshTokenExp,
		refreshTokenRepo: refreshTokenRepo,
	}
}

// GenerateAccessToken generates a new access token for the given user ID
// using the private key and the specified expiration time.
// The token is signed with the RS256 algorithm and includes the user ID and issuer as claims.
func (s *JWTTokenService) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    s.issuer,
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExp)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.private)
}

// VerifyAccessToken verifies the given access token string using the public key.
// It checks the token's signature, expiration time, and issuer.
// If the token is valid, it returns the user ID from the token claims.
func (s *JWTTokenService) VerifyAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.public, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Issuer != s.issuer {
			return "", fmt.Errorf("invalid issuer: %s", claims.Issuer)
		}
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return "", fmt.Errorf("token expired")
		}
		return claims.Subject, nil
	}

	return "", fmt.Errorf("invalid token")
}

// GenerateRefreshToken generates a new refresh token for the given user ID.
// It generates a random string of 64 characters and stores it in the refresh token repository
// with an expiration time of 7 days. It returns the generated token and its expiration time.
func (s *JWTTokenService) GenerateRefreshToken(ctx context.Context, userID string) (string, time.Time, error) {
	token, err := generateRandomString(64)

	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	expiresAt := time.Now().Add(24 * 7 * time.Hour) // 7 days

	// Store the refresh token in the repository with an expiration time
	if err := s.refreshTokenRepo.Store(ctx, userID, token, expiresAt); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to store refresh token: %w", err)
	}
	return token, expiresAt, nil
}

// VerifyRefreshToken verifies the given refresh token string.
// It checks if the token exists in the refresh token repository and if it has not expired.
func (s *JWTTokenService) VerifyRefreshToken(ctx context.Context, token string) (string, error) {
	userID, expiredAt, err := s.refreshTokenRepo.Get(ctx, token)
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}

	if time.Now().After(expiredAt) {
		err := s.refreshTokenRepo.Revoke(ctx, token)
		if err != nil {
			return "", fmt.Errorf("failed to revoke refresh token: %w", err)
		}
		return "", fmt.Errorf("refresh token expired")
	}
	return userID, nil
}

// RevokeRefreshToken invalidates a refresh token
func (s *JWTTokenService) RevokeRefreshToken(ctx context.Context, token string) error {
	return s.refreshTokenRepo.Revoke(ctx, token)
}

// RotateRefreshToken generates a new refresh token for the given old token.
// It verifies the old token, revokes it, and generates a new one.
// It returns the new token and its expiration time.
func (s *JWTTokenService) RotateRefreshToken(ctx context.Context, oldToken string) (string, time.Time, error) {
	userID, err := s.VerifyRefreshToken(ctx, oldToken)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to verify refresh token: %w", err)
	}

	err = s.RevokeRefreshToken(ctx, oldToken)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	newToken, expiredAt, err := s.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	return newToken, expiredAt, nil
}
