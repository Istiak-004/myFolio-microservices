package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/istiak-004/myFolio-microservices/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

type Client struct {
	// DB is the database connection
	db *sqlx.DB
	// Config is the database configuration
	config *Config
	// Logger is the logger
	logger *logger.Logger
}

var (
	instance *Client
	once     sync.Once
)

// Config holds database configuration
type Config struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// New creates a new database client
func New(config *Config, logger *logger.Logger) (*Client, error) {
	var initErr error
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode)

		db, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			initErr = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		// Configure connection pool
		db.SetMaxOpenConns(config.MaxOpenConns)
		db.SetMaxIdleConns(config.MaxIdleConns)
		db.SetConnMaxLifetime(config.ConnMaxLifetime)

		// Set connection timeout
		ctx, cancel := context.WithTimeout(context.Background(), config.ConnMaxLifetime)
		defer cancel()
		// Ping the database to verify connection
		if err := db.PingContext(ctx); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}
		instance = &Client{
			db:     db,
			config: config,
			logger: logger,
		}
		logger.Info("Database connection established")
		logger.WithComponent("database").Info("Database connection established")
	})

	return instance, initErr
}

// GetDB returns the underlying sqlx.DB instance
func (c *Client) GetDB() *sqlx.DB {
	return c.db
}

// Close gracefully closes the database connection
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
