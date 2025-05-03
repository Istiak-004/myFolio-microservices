package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// loadConfig loads configuration from file and environment
func loadConfig(v *viper.Viper) error {
	// Try to read from config file
	if err := v.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist, we'll use defaults + env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Override with environment variables if they exist
	for _, key := range v.AllKeys() {
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		if _, ok := os.LookupEnv("MYAPP_" + envKey); ok {
			if err := v.BindEnv(key, "MYAPP_"+envKey); err != nil {
				return fmt.Errorf("failed to bind env var for key %s: %w", key, err)
			}
		}
	}

	return nil
}

// bindEnvVars sets up explicit environment variable bindings
func bindEnvVars(v *viper.Viper) {
	// Application
	v.BindEnv("app.environment", "MYAPP_ENV")

	// HTTP
	v.BindEnv("http.port", "MYAPP_HTTP_PORT")
	v.BindEnv("http.timeout", "MYAPP_HTTP_TIMEOUT")

	// Database
	v.BindEnv("database.host", "MYAPP_DB_HOST")
	v.BindEnv("database.port", "MYAPP_DB_PORT")
	v.BindEnv("database.user", "MYAPP_DB_USER")
	v.BindEnv("database.password", "MYAPP_DB_PASSWORD")
	v.BindEnv("database.name", "MYAPP_DB_NAME")
	v.BindEnv("database.sslmode", "MYAPP_DB_SSLMODE")

	// Logging
	v.BindEnv("log.level", "MYAPP_LOG_LEVEL")
	v.BindEnv("log.format", "MYAPP_LOG_FORMAT")
}
