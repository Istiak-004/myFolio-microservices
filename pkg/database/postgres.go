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
	config *SQLConfig
	// Logger is the logger
	logger *logger.Logger
}

var (
	instance *Client
	once     sync.Once
)

type DBConfigRepo interface {
	GetConfig() *SQLConfig
}

// Config holds database configuration
type SQLConfig struct {
	Host            string        `mapstructure:"DB_HOST"`
	Port            int           `mapstructure:"DB_PORT"`
	User            string        `mapstructure:"DB_USER"`
	Password        string        `mapstructure:"DB_PASSWORD"`
	Name            string        `mapstructure:"DB_NAME"`
	SSLMode         string        `mapstructure:"DB_SSL_MODE"`
	MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

// New creates a new database client
func New[T DBConfigRepo](config T, logger *logger.Logger) (*Client, error) {
	var initErr error
	once.Do(func() {
		dbConfig := config.GetConfig()
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode)

		db, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			initErr = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		// dbConfigure connection pool
		db.SetMaxOpenConns(dbConfig.MaxOpenConns)
		db.SetMaxIdleConns(dbConfig.MaxIdleConns)
		db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

		// Set connection timeout
		ctx, cancel := context.WithTimeout(context.Background(), dbConfig.ConnMaxLifetime)
		defer cancel()
		// Ping the database to verify connection
		if err := db.PingContext(ctx); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}
		instance = &Client{
			db:     db,
			config: dbConfig,
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
