package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load loads the base and service-specific config into the target object.
// Usage: Load("auth", &AuthConfig{})
// LoadConfig loads config from a specific directory path
func LoadConfig[T Config](prefix, configPath string, target T) (Base, T, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	v.AutomaticEnv()
	v.SetEnvPrefix(strings.ToUpper(prefix))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Safe defaults
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("ENV", "development")
	v.SetDefault("APP_NAME", prefix)

	// Optional config file
	if err := v.ReadInConfig(); err != nil {
		return Base{}, target, fmt.Errorf("failed to read config file: %w", err)
	}

	var base Base
	if err := v.Unmarshal(&base); err != nil {
		return Base{}, target, fmt.Errorf("base config load failed: %w", err)
	}

	if err := v.Unmarshal(&target); err != nil {
		return Base{}, target, fmt.Errorf("service config load failed: %w", err)
	}

	if err := target.Validate(); err != nil {
		return Base{}, target, fmt.Errorf("service config invalid: %w", err)
	}

	return base, target, nil
}
