package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	errInvalidToken = errors.New("invalid or expired token")
)

type TokenManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenManager(privateKeyPath, publicKeyPath, issuer string, accessTTL, refreshTTL time.Duration) (*TokenManager, error) {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}
	publicKey, err := loadPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}
	return &TokenManager{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

func (tm *TokenManager) GenerateAccessToken(userID, role string) (string, string, error) {
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

func (tm *TokenManager) GenerateRefreshToken(userID string) (string, string, error) {
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

func (tm *TokenManager) ParseToken(tokenStr string) (*CustomClaims, error) {
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

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key format")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := ioutil.ReadFile(path)
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
