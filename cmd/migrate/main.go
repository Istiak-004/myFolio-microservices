package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"strconv"

	"github.com/istiak-004/myFolio-microservices/pkg/database"
	"github.com/istiak-004/myFolio-microservices/pkg/logger"
	_ "github.com/lib/pq"
)

func main() {
	// Parse command line flags
	up := flag.Bool("up", false, "Apply all up migrations")
	down := flag.Bool("down", false, "Apply all down migrations")
	create := flag.String("create", "", "Create new migration with given name")
	_ = flag.String("config", "./configs", "Path to config directory")
	service := flag.String("service", "", "Service name (auth, notification, etc.)")
	flag.Parse()

	if *service == "" {
		log.Fatal("Service name must be provided with -service flag")
	}

	// Initialize logger
	logger := logger.NewLogger("Migrate Service")
	logger.Info("Starting migration service", logger.String("service", *service))

	// Load database config
	dbConfig := &database.Config{
		Host:            os.Getenv("DB_HOST"),
		Port:            mustAtoi(os.Getenv("DB_PORT")),
		User:            os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		Name:            os.Getenv("DB_NAME"),
		SSLMode:         os.Getenv("DB_SSLMODE"),
		MaxOpenConns:    mustAtoi(os.Getenv("DB_MAX_OPEN_CONNS")),
		MaxIdleConns:    mustAtoi(os.Getenv("DB_MAX_IDLE_CONNS")),
		ConnMaxLifetime: mustParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME")),
	}

	// Create database client
	dbClient, err := database.New(dbConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create database client", logger.ErrorFields(err))
	}
	defer dbClient.Close()

	// Create migrator
	migrationsPath := fmt.Sprintf("./services/%s/migrations", *service)
	migrator, err := database.NewMigrator(dbClient.GetDB(), migrationsPath)
	if err != nil {
		logger.Fatal("Failed to create migrator")
	}

	// Execute command
	switch {
	case *up:
		if err := migrator.Up(); err != nil {
			logger.Fatal("Failed to apply migrations", logger.ErrorFields(err))
		}
		logger.Info("Migrations applied successfully")
	case *down:
		if err := migrator.Down(); err != nil {
			logger.Fatal("Failed to rollback migrations", logger.ErrorFields(err))
		}
		logger.Info("Migrations rolled back successfully")
	case *create != "":
		if err := migrator.CreateMigration(*create); err != nil {
			logger.Fatal("Failed to create migration", logger.ErrorFields(err))
		}
		logger.Info("Migration files created")
	default:
		logger.Fatal("No command specified (use -up, -down, or -create)")
	}
}

// mustParseDuration parses a duration string and panics if there's an error.
func mustParseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("Failed to parse duration: %v", err)
	}
	return duration
}

// mustAtoi converts a string to an integer and panics if there's an error.
func mustAtoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %v", err)
	}
	return value
}
