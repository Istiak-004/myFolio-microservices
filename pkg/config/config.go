package config

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/istiak-004/myFolio-microservices/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// singleton instance of viper
// to ensure that the configuration is loaded only once
// and is thread-safe
// using sync.Once to ensure that the instance is created only once
// and is safe for concurrent use
var (
	configInstance *Config
	once           sync.Once
	configLogger   *logger.Logger // Package-level logger instance
)

// Config represents the application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	HTTP     HTTPConfig
	Log      LogConfig
}

// Init initializes the configuration package with logging
func Init(serviceName string) (*Config, error) {
	var initErr error
	once.Do(func() {
		// Initialize logger first
		configLogger = logger.NewLogger(serviceName)
		configLogger.Info("Initializing configuration...")

		configInstance = &Config{}

		// Set up Viper
		v := viper.New()
		if err := setupViper(v); err != nil {
			configLogger.WithError(err).Error("Failed to setup Viper configuration")
			initErr = fmt.Errorf("failed to setup viper: %w", err)
			return
		}

		// Load configuration
		if err := loadConfig(v); err != nil {
			configLogger.WithError(err).Error("Failed to load configuration")
			initErr = fmt.Errorf("failed to load config: %w", err)
			return
		}

		// Unmarshal configuration
		if err := v.Unmarshal(configInstance); err != nil {
			configLogger.WithError(err).
				WithField("component", "config-unmarshal").
				Error("Failed to unmarshal configuration")
			initErr = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}

		configLogger.WithFields(logrus.Fields{
			"environment": configInstance.App.Environment,
			"version":     configInstance.App.Version,
		}).Info("Configuration initialized successfully")
	})

	return configInstance, initErr
}

// Get returns the configuration instance (must call Init first)
func Get() *Config {
	if configInstance == nil {
		panic("config not initialized - call Init() first")
	}
	return configInstance
}

// setupViper configures Viper with default settings
func setupViper(v *viper.Viper) {
	// Set default values
	setDefaults(v)

	// Configuration file name (without extension)
	v.SetConfigName("config")

	// Configuration type
	v.SetConfigType("yaml")

	// Paths to look for the config file
	v.AddConfigPath(".")                           // Current directory
	v.AddConfigPath(filepath.Join("..", "config")) // Parent directory's config folder
	v.AddConfigPath("/etc/myapp/")                 // System config directory
	v.AddConfigPath("$HOME/.myapp")                // User config directory

	// Enable environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("MYAPP") // Environment variables will be prefixed with MYAPP_

	// Configure environment variable bindings
	bindEnvVars(v)
}
