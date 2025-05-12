package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/istiak-004/myFolio-microservices/pkg/logger"
	"github.com/redis/go-redis/v9"
)

// RedisConfig holds the Redis configuration
type RedisConfig struct {
	Host         string        `mapstructure:"REDIS_HOST"`
	Password     string        `mapstructure:"REDIS_PASSWORD"` // Use Secrets Manager in prod
	DialTimeout  time.Duration `mapstructure:"REDIAL_DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"REDIS_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"REDIS_WRITE_TIMEOUT"`
	Port         int           `mapstructure:"REDIS_PORT"`
	DB           int           `mapstructure:"REDIS_DB"`
	PoolSize     int           `mapstructure:"REDIS_POOL_SIZE"`
}

type HasRedisCionfig interface {
	GetRedisConfig() *RedisConfig
}

type RedisClient struct {
	client *redis.Client
	config *RedisConfig
	logger *logger.Logger
}

var (
	redisInstance *RedisClient
	redisOnce     sync.Once
)

// RedisManager interface for testability and flexibility
type RedisManager interface {
	GetClient() *redis.Client
	Close() error
}

// NewRedisClient creates a new Redis client
func NewRedisClient[T HasRedisCionfig](config T, logger *logger.Logger) (*RedisClient, error) {
	var initErr error
	redisOnce.Do(func() {
		redisConfig := config.GetRedisConfig()
		addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
		rdb := redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     redisConfig.Password,
			DB:           redisConfig.DB,
			DialTimeout:  redisConfig.DialTimeout,
			ReadTimeout:  redisConfig.ReadTimeout,
			WriteTimeout: redisConfig.WriteTimeout,
			PoolSize:     redisConfig.PoolSize,
		})

		ctx, err := context.WithTimeout(context.Background(), redisConfig.DialTimeout)
		if err != nil {
			initErr = fmt.Errorf("failed to create context: %w", err)
			return
		}

		if err := rdb.Ping(ctx).Err(); err != nil {
			initErr = fmt.Errorf("redis ping failed: %w", err)
			return
		}

		redisInstance = &RedisClient{
			client: rdb,
			config: redisConfig,
			logger: logger,
		}

		logger.WithComponent("redis").Info("Redis connection established")
	})

	return redisInstance, initErr
}

// GetClient returns the underlying Redis client
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// Close gracefully closes the Redis client
func (r *RedisClient) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
