// config/config.go
package config

import (
	"time"

	"github.com/istiak-004/myFolio-microservices/pkg/logger"
)

type AuthConfig struct {
	AppName string `mapstructure:"APP_NAME"`
	Env     string `mapstructure:"ENV"`
	Port    string `mapstructure:"PORT"`

	DB     DBConfig
	Redis  RedisConfig
	SMTP   SMTPConfig
	JWT    JWTConfig
	OAuth  OAuthConfig
	Logger *logger.Logger
}

type DBConfig struct {
	Host            string `mapstructure:"DB_HOST"`
	Port            int    `mapstructure:"DB_PORT"`
	User            string `mapstructure:"DB_USER"`
	Password        string `mapstructure:"DB_PASSWORD"`
	Name            string `mapstructure:"DB_NAME"`
	SSLMode         string `mapstructure:"DB_SSL_MODE"`
	MaxOpenConns    int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime string `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"ADDR"`
	Password string `mapstructure:"PASSWORD"`
	DB       int    `mapstructure:"DB"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"SMTP_HOST"`
	Port     int    `mapstructure:"SMTP_PORT"`
	Username string `mapstructure:"SMTP_USERNAME"`
	Password string `mapstructure:"SMTP_PASSWORD"`
	From     string `mapstructure:"SMTP_FROM"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"JWT_SECRET"`
	Expiry time.Duration `mapstructure:"JWT_EXPIRY"`
}

type OAuthConfig struct {
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GithubClientID     string `mapstructure:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET"`
	RedirectURL        string `mapstructure:"REDIRECT_URL"`
}
