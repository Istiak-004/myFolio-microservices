package config

import "github.com/spf13/viper"

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Application defaults
	v.SetDefault("app.name", "myFolio-microservices")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.version", "1.0.0")

	// HTTP server defaults
	v.SetDefault("http.port", 8080)
	v.SetDefault("http.timeout", 30)
	v.SetDefault("http.shutdown_timeout", 10)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "myapp")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")

	// Logging defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")
}
