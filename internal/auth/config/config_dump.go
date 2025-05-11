package config

import "time"

type ConfigT struct {
	Server struct {
		Port         string        `env:"PORT" envDefault:"8080"`
		ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"30s"`
		WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"30s"`
	}
	Database struct {
		URL           string `env:"DB_URL" envDefault:"postgres://user:pass@localhost:5432/auth?sslmode=disable"`
		MigrationsDir string `env:"MIGRATIONS_DIR" envDefault:"migrations"`
	}
	Redis struct {
		URL string `env:"REDIS_URL" envDefault:"redis://localhost:6379/0"`
	}
	JWT struct {
		PrivateKeyPath string        `env:"JWT_PRIVATE_KEY" envDefault:"./keys/private.pem"`
		PublicKeyPath  string        `env:"JWT_PUBLIC_KEY" envDefault:"./keys/public.pem"`
		AccessExpiry   time.Duration `env:"JWT_ACCESS_EXPIRY" envDefault:"15m"`
		RefreshExpiry  time.Duration `env:"JWT_REFRESH_EXPIRY" envDefault:"168h"` // 7 days
	}
	SMTP struct {
		Host     string `env:"SMTP_HOST"`
		Port     int    `env:"SMTP_PORT" envDefault:"587"`
		Username string `env:"SMTP_USERNAME"`
		Password string `env:"SMTP_PASSWORD"`
		From     string `env:"SMTP_FROM" envDefault:"noreply@example.com"`
	}
	OAuth struct {
		Google struct {
			ClientID     string `env:"OAUTH_GOOGLE_ID"`
			ClientSecret string `env:"OAUTH_GOOGLE_SECRET"`
			RedirectURL  string `env:"OAUTH_GOOGLE_REDIRECT"`
		}
		GitHub struct {
			ClientID     string `env:"OAUTH_GITHUB_ID"`
			ClientSecret string `env:"OAUTH_GITHUB_SECRET"`
			RedirectURL  string `env:"OAUTH_GITHUB_REDIRECT"`
		}
	}
	App struct {
		BaseURL          string `env:"APP_BASE_URL" envDefault:"http://localhost:8080"`
		VerificationPath string `env:"VERIFICATION_PATH" envDefault:"/verify-email"`
	}
}
