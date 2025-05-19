package config

type Redis struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	DB       int    `mapstructure:"db" validate:"required"`
	PoolSize int    `mapstructure:"pool_size" validate:"required"`
}

type SMTP struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	From     string `mapstructure:"from" validate:"required"`
}

type AuthConfig struct {
	App          AppConfig      `mapstructure:"app" validate:"required"`
	Database     DatabaseConfig `mapstructure:"database" validate:"required,dive"`
	HTTP         HTTPConfig     `mapstructure:"http" validate:"required,dive"`
	Log          LogConfig      `mapstructure:"log" validate:"required,dive"`
	JWTIssuer    string         `mapstructure:"jwt_issuer" validate:"required"`
	JWTExpiry    string         `mapstructure:"jwt_expiry" validate:"required"`
	CookieDomain string         `mapstructure:"cookie_domain" validate:"required"`
}

type AppConfig struct {
	Name        string `mapstructure:"name" validate:"required"`
	Environment string `mapstructure:"environment" validate:"required"`
	Version     string `mapstructure:"version" validate:"required"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required"`
	User            string `mapstructure:"user" validate:"required"`
	Password        string `mapstructure:"password" validate:"required"`
	Name            string `mapstructure:"name" validate:"required"`
	SSLMode         string `mapstructure:"sslmode" validate:"required"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" validate:"required"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime" validate:"required"`
}

// HTTPConfig represents HTTP server configuration
type HTTPConfig struct {
	Port            int `mapstructure:"port" validate:"required"`
	Timeout         int `mapstructure:"timeout" validate:"required"`
	ShutdownTimeout int `mapstructure:"shutdown_timeout" validate:"required"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level" validate:"required"`
	Format string `mapstructure:"format" validate:"required"`
	Output string `mapstructure:"output" validate:"required"`
}
