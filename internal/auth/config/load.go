package config

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func (a *AuthConfig) Validate() error {
	validator := validator.New()
	if err := validator.Struct(a); err != nil {
		return err
	}
	return nil
}

// Implement the database.Config interface
func (c *DatabaseConfig) GetHost() string {
	return c.Host
}

func (c *DatabaseConfig) GetPort() int {
	return c.Port
}

func (c *DatabaseConfig) GetUser() string {
	return c.User
}

func (c *DatabaseConfig) GetPassword() string {
	return c.Password
}

func (c *DatabaseConfig) GetName() string {
	return c.Name
}

func (c *DatabaseConfig) GetSSLMode() string {
	return c.SSLMode
}

func (c *DatabaseConfig) GetMaxOpenConns() int {
	return c.MaxOpenConns
}

func (c *DatabaseConfig) GetMaxIdleConns() int {
	return c.MaxIdleConns
}

func (c *DatabaseConfig) GetConnMaxLifetime() time.Duration {
	duration, err := time.ParseDuration(c.ConnMaxLifetime)
	if err != nil {
		// Return a default value or handle the error as appropriate
		return 30 * time.Minute
	}
	return duration
}
