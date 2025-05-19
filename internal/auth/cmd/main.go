package main

import (
	"fmt"

	authconfig "github.com/istiak-004/myFolio-microservices/auth/config"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
	"github.com/istiak-004/myFolio-microservices/pkg/config"
	"github.com/istiak-004/myFolio-microservices/pkg/logger"
)

func main() {
	log := logger.NewLogger("Auth")
	log.Info("Starting Auth Service", log.String("service", "auth"))

	// Load configuration
	log.Info("Loading configuration")
	// Load the configuration using the config package
	var authConfig authconfig.AuthConfig

	base, config, err := config.LoadConfig("auth", "../config", &authConfig)
	if err != nil {
		log.WithError(err).Error("Failed to initialize configuration")
		return
	}
	log.WithFields(map[string]interface{}{
		"base":   base,
		"config": config,
	}).Info("Configuration loaded successfully")

	u := models.User{}
	fmt.Println(u)
}
