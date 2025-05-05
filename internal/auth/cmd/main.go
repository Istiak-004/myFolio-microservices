package main

import (
	"github.com/istiak-004/myFolio-microservices/pkg/config"
	"github.com/istiak-004/myFolio-microservices/pkg/logger"
)

func main() {
	log:=logger.NewLogger("Auth Service")
	log.Info("Starting Auth Service",log.String("service", "auth"))
	config, err := config.Init("Auth Service")
	if err != nil {
		log.WithError(err).Error("Failed to initialize configuration")
		return
	}
	log.WithFields(map[string]interface{}{
		"config": config,
	}).Info("Configuration loaded successfully")
}
