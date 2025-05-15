package token

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	errInvalidToken = errors.New("invalid or expired token")
)

// TokenManager is responsible for generating and verifying JWT tokens
type JWTTokenManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type CustomClaims struct {
	UserID               string `json:"user_id"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // embedded standard claims
}

// NewTokenManager creates a new TokenManager with the given parameters
func NewTokenManager(privateKeyPath, publicKeyPath, issuer string, accessTTL, refreshTTL time.Duration) (*JWTTokenManager, error) {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}
	publicKey, err := loadPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}
	return &JWTTokenManager{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

// GenerateAccessToken generates a new access token for the given user ID and role
func (tm *JWTTokenManager) GenerateAccessToken(userID, role string) (string, string, error) {
	jti := uuid.New().String()
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   userID,
			Issuer:    tm.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(tm.privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, jti, nil
}

// GenerateRefreshToken generates a new refresh token for the given user ID
func (tm *JWTTokenManager) GenerateRefreshToken(userID string) (string, string, error) {
	jti := uuid.New().String()
	claims := jwt.RegisteredClaims{
		ID:        jti,
		Subject:   userID,
		Issuer:    tm.issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.refreshTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(tm.privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return signed, jti, nil
}

// ParseAccessToken parses and validates the access token
func (tm *JWTTokenManager) ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errInvalidToken
	}
	return claims, nil
}

// VerifyAccessToken verifies the access token and returns the claims
func (tm *JWTTokenManager) VerifyAccessToken(ctx context.Context, tokenStr string) (*CustomClaims, error) {
	return tm.ParseToken(tokenStr)
}

// loadPrivateKey loads the private key from the given path
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key format")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// loadPublicKey loads the public key from the given path
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key format")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return key, nil
}
