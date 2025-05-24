package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	authconfig "github.com/istiak-004/myFolio-microservices/auth/config"
	"github.com/istiak-004/myFolio-microservices/pkg/config"
	"github.com/istiak-004/myFolio-microservices/pkg/database"
	"github.com/istiak-004/myFolio-microservices/pkg/logger"

	// "github.com/istiak-004/myFolio-microservices/pkg/server"
	"go.uber.org/zap"
)

const (
	serviceName     = "auth"
	shutdownTimeout = 15 * time.Second
)

func main() {
	// Initialize context with cancellation
	_, cancel := context.WithCancel(context.Background())
	
	defer cancel()

	// Initialize logger
	log := logger.NewLogger(serviceName)
	defer log.Sync() // Ensure logs are flushed

	// Setup signal handling for graceful shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	// Load configuration
	log.Info("Loading configuration...")
	var authConfig authconfig.AuthConfig
	baseConfig, appConfig, err := config.LoadConfig(serviceName, "../config", &authConfig)
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	log.Info("Configuration loaded successfully",
		zap.Any("base_config", baseConfig),
		zap.Any("app_config", appConfig),
	)

	// Initialize database
	log.Info("Initializing database connection...")
	dbClient, err := database.NewDB(&appConfig.Database, log)
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer func() {
		log.Info("Closing database connection...")
		if err := dbClient.Close(); err != nil {
			log.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	log.Info("Database connection established successfully")

	// // Initialize services
	// // Example: authService := auth.NewService(dbClient, log)

	// // Initialize HTTP server
	// log.Info("Starting HTTP server...")
	// srv := server.NewHTTPServer(appConfig.Server, log /*, authService */)

	// go func() {
	// 	if err := srv.Start(); err != nil {
	// 		log.Error("HTTP server error", zap.Error(err))
	// 		shutdownCh <- syscall.SIGTERM
	// 	}
	// }()

	// log.Info("Service started successfully",
	// 	zap.String("service", serviceName),
	// 	zap.String("version", baseConfig.Version),
	// )

	// // Wait for shutdown signal
	// sig := <-shutdownCh
	// log.Info("Received shutdown signal", zap.String("signal", sig.String()))

	// // Graceful shutdown
	// gracefulShutdown(ctx, log, srv, dbClient)
}

// func gracefulShutdown(
// 	ctx context.Context,
// 	log *zap.Logger,
// 	srv *server.HTTPServer,
// 	dbClient database.DatabaseManager,
// ) {
// 	// Create a deadline for shutdown operations
// 	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
// 	defer cancel()

// 	log.Info("Initiating graceful shutdown...")

// 	// Shutdown HTTP server
// 	if srv != nil {
// 		if err := srv.Shutdown(shutdownCtx); err != nil {
// 			log.Error("HTTP server shutdown error", zap.Error(err))
// 		} else {
// 			log.Info("HTTP server stopped gracefully")
// 		}
// 	}

// 	// Close database connection
// 	if dbClient != nil {
// 		if err := dbClient.Close(); err != nil {
// 			log.Error("Database connection closure error", zap.Error(err))
// 		} else {
// 			log.Info("Database connection closed gracefully")
// 		}
// 	}

// 	log.Info("Service shutdown completed")
// }
