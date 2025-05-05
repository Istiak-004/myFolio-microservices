package token

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepo struct {
	redisClient *redis.Client
	prefix      string
}

func NewRefreshTokenRepo(
	redisClient *redis.Client,
	ttl time.Duration,
	prefix string,
) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		redisClient: redisClient,
		prefix:      prefix,
	}
}

// store the refresh token in Redis with an expiration time
// and add it to the user's set of tokens
func (r *RefreshTokenRepo) Store(ctx context.Context, userID, token string, expiresAt time.Time) error {

	tokenKeyString := r.prefix + "token:" + token
	// Store the refresh token in Redis with an expiration time
	if err := r.redisClient.Set(ctx, tokenKeyString, userID, time.Until(expiresAt)).Err(); err != nil {
		return err
	}

	userKeyString := r.prefix + "user_tokens:" + userID

	// Add the token to the user's set in Redis
	// This will automatically remove the token from the set when it expires
	if err := r.redisClient.ZAdd(ctx, userKeyString, redis.Z{
		Score:  float64(expiresAt.Unix()),
		Member: token,
	}).Err(); err != nil {
		// If ZAdd fails, we should delete the token from Redis
		r.redisClient.Del(ctx, tokenKeyString)
		return fmt.Errorf("failed to add token to user set: %w", err)
	}

	// Set the expiration time for the user ID key
	// This will remove the user ID key if it has no tokens left
	r.redisClient.Expire(ctx, userKeyString, time.Until(expiresAt))
	return nil
}

// get the user ID and expiration time for a given refresh token
// If the token is not found, return an error
// If the token is found, return the user ID and expiration time
// If the token is found but has expired, return an error
func (r *RefreshTokenRepo) Get(ctx context.Context, token string) (string, time.Time, error) {
	tokenKeyString := r.prefix + "token:" + token

	// Use a pipeline to get the user ID and TTL in one round trip
	// This is more efficient than getting them separately
	pipe := r.redisClient.Pipeline()
	userIDCmd := pipe.Get(ctx, tokenKeyString)
	ttlCmd := pipe.TTL(ctx, tokenKeyString)
	_, err := pipe.Exec(ctx)
	if err != nil {
		if err == redis.Nil {
			return "", time.Time{}, fmt.Errorf("token not found")
		}
		return "", time.Time{}, fmt.Errorf("failed to get token: %w", err)
	}

	userID, err := userIDCmd.Result()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get user ID: %w", err)
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get TTL: %w", err)
	}
	return userID, time.Now().Add(ttl), nil
}

// revoke a refresh token by deleting it from Redis
// This will also remove it from the user's set of tokens
func (r *RefreshTokenRepo) Revoke(ctx context.Context, token string) error {
	userID, _, err := r.Get(ctx, token)
	if err != nil {
		if err.Error() == "token not found" {
			return nil // Token already revoked
		}
		return fmt.Errorf("failed to get token: %w", err)
	}

	// delete the token from the user's set
	tokenKeyString := r.prefix + "token:" + token
	if err := r.redisClient.Del(ctx, tokenKeyString).Err(); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	// delete the token from the user's set
	userKeyString := r.prefix + "user_tokens:" + userID
	if err := r.redisClient.ZRem(ctx, userKeyString, token).Err(); err != nil {
		return fmt.Errorf("failed to remove token from user set: %w", err)
	}

	return nil
}

// revoke all refresh tokens for a given user ID
// This will delete all tokens from Redis and remove the user's set of tokens
func (r *RefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID string) error {
	userKeyString := r.prefix + "user_tokens:" + userID

	// Get all tokens for the user
	tokens, err := r.redisClient.ZRange(ctx, userKeyString, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get user tokens: %w", err)
	}

	// Delete each token from Redis
	for _, token := range tokens {
		tokenKeyString := r.prefix + "token:" + token
		if err := r.redisClient.Del(ctx, tokenKeyString).Err(); err != nil {
			return fmt.Errorf("failed to delete token: %w", err)
		}
	}

	// Delete the user's set of tokens
	if err := r.redisClient.Del(ctx, userKeyString).Err(); err != nil {
		return fmt.Errorf("failed to delete user tokens: %w", err)
	}

	return nil
}
