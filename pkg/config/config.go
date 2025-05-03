package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// singleton instance of viper
// to ensure that the configuration is loaded only once
// and is thread-safe
// using sync.Once to ensure that the instance is created only once
// and is safe for concurrent use
var (
	instance *viper.Viper
	once     sync.Once
)

// GetConfig returns the singleton instance of viper
func GetConfig() *viper.Viper {
	once.Do(func() {
		instance = viper.New()
		instance.SetConfigName("config")
		instance.SetConfigType("yaml")
		instance.AddConfigPath(".")
		instance.AddConfigPath("./config")
		instance.AutomaticEnv()

		// Set default values for configuration
		instance.SetDefault("server.port", 8080)
		instance.SetDefault("database.max_conns", 20)

		if err := instance.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				panic(fmt.Errorf("Error reading config file: %s", err))
			}
		}
	})
	return instance
}
