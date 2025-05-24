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
	config Config
	// Logger is the logger
	logger *logger.Logger
}

var (
	instance *Client
	once     sync.Once
)

// Config defines the minimum required database configuration interface
type Config interface {
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetName() string
	GetSSLMode() string
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxLifetime() time.Duration
}

type DatabaseManager interface {
	GetDB() *sqlx.DB
	Close() error
}

// New creates a new database client
// It uses a singleton pattern to ensure only one instance of the database client is created
// and reused throughout the application.
func NewDB(config Config, logger *logger.Logger) (*Client, error) {
	var initErr error
	once.Do(func() {

		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.GetHost(), config.GetPort(), config.GetUser(),
			config.GetPassword(), config.GetName(), config.GetSSLMode())

		db, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			initErr = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		// dbConfigure connection pool
		db.SetMaxOpenConns(config.GetMaxOpenConns())
		db.SetMaxIdleConns(config.GetMaxIdleConns())
		db.SetConnMaxLifetime(config.GetConnMaxLifetime())

		// Set connection timeout
		ctx, cancel := context.WithTimeout(context.Background(), config.GetConnMaxLifetime())
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
		// logger.Info("Database connection established")
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
