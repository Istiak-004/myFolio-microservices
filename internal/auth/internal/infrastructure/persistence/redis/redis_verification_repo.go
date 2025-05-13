package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/istiak-004/myFolio-microservices/pkg/database"
	"github.com/redis/go-redis/v9"
)

const (
	refreshTokenPrefix = "auth:refresh_token:"
	usedJtiPrefix      = "auth:used_jti:"
	tokenTTL           = 30 * 24 * time.Hour // 30 days
	usedJTITTL         = 60 * time.Minute    // prevent reuse for 1 hour
)

type TokenRepository struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *database.RedisClient) *TokenRepository {
	redisClient := rdb.GetClient()
	return &TokenRepository{rdb: redisClient}
}

// hashToken hashes the token using SHA-256 and returns the hex string.
// This is used to create a unique identifier for the token in Redis.
// It is important to use a secure hashing algorithm to prevent token collisions.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// StoreRefreshToken stores the refresh token in Redis with an expiration time.
func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	hashed := hashToken(token)
	jti := fmt.Sprintf("%s%s", refreshTokenPrefix, hashed)
	return r.rdb.Set(ctx, jti, userID, time.Until(expiresAt)).Err()
}

// VerifyRefreshToken checks if the refresh token is valid and returns the associated user ID.
// It uses the hashed token to look up the user ID in Redis.
func (r *TokenRepository) VerifyRefreshToken(ctx context.Context, token string) (string, error) {
	hashed := hashToken(token)
	jti := fmt.Sprintf("%s%s", refreshTokenPrefix, hashed)
	userID, err := r.rdb.Get(ctx, jti).Result()
	if err == redis.Nil {
		return "", errors.New("token not found or expired")
	} else if err != nil {
		return "", err
	}
	return userID, nil
}

// RevokeRefreshToken removes the refresh token from Redis, effectively invalidating it.
func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	hashed := hashToken(token)
	jti := fmt.Sprintf("%s%s", refreshTokenPrefix, hashed)
	return r.rdb.Del(ctx, jti).Err()
}

// MarkJTIAssigned marks a JTI (JWT ID) as used in Redis to prevent replay attacks.
// It sets a key with the JTI and a TTL to ensure it is not reused.
// This is useful for JWTs where you want to ensure that a token can only be used once.
func (r *TokenRepository) MarkJTIAssigned(ctx context.Context, jti string) error {
	key := fmt.Sprintf("%s%s", usedJtiPrefix, jti)
	return r.rdb.Set(ctx, key, true, usedJTITTL).Err()
}

// IsJTIUsed checks if a JTI (JWT ID) has already been used.
// It checks if the key exists in Redis and returns true if it does, false otherwise.
// This is useful for preventing replay attacks where a token could be reused.
func (r *TokenRepository) IsJTIUsed(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("%s%s", usedJtiPrefix, jti)
	exists, err := r.rdb.Exists(ctx, key).Result()
	return exists > 0, err
}

// RotateRefreshToken revokes the old refresh token and stores a new one.
// It returns the new token and any error encountered during the process.
// This is useful for implementing refresh token rotation, where a new token is issued
// each time the old one is used. This helps to prevent replay attacks and ensures that
// the refresh token is always fresh.
func (r *TokenRepository) RotateRefreshToken(ctx context.Context, oldToken string, newToken string, userID string) (string, error) {
	// Revoke old
	_ = r.RevokeRefreshToken(ctx, oldToken)

	// Store new
	expiresAt := time.Now().Add(tokenTTL)
	err := r.StoreRefreshToken(ctx, userID, newToken, expiresAt)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

func (r *TokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	// Optional: Use SCAN to find all keys by userID if you structure keys differently
	// For now, this is a no-op (requires more tracking)
	return nil // OR: implement tracking jtis per user
}
