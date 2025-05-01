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

func (r *RefreshTokenRepo) Get(ctx context.Context, token string) (string, time.Time, error) {
	tokenKeyString := r.prefix + "token:" + token

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
