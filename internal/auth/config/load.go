package config

import (
	"fmt"

	"go.uber.org/zap"
)

func (a *AuthConfig) Validate() error {

	missingFields := []string{}
	// Validate required fields
	if a.AppName == "" {
		missingFields = append(missingFields, "APP_NAME")
	}
	if a.Env == "" {
		missingFields = append(missingFields, "ENV")
	}
	if a.Port == "" {
		missingFields = append(missingFields, "PORT")
	}
	if a.DB.Host == "" {
		missingFields = append(missingFields, "DB_HOST")
	}
	if a.DB.Port == 0 {
		missingFields = append(missingFields, "DB_PORT")
	}
	if a.DB.User == "" {
		missingFields = append(missingFields, "DB_USER")
	}
	if a.DB.Name == "" {
		missingFields = append(missingFields, "DB_NAME")
	}
	if a.Redis.Addr == "" {
		missingFields = append(missingFields, "REDIS_ADDR")
	}
	if a.SMTP.Host == "" {
		missingFields = append(missingFields, "SMTP_HOST")
	}
	if a.SMTP.Username == "" {
		missingFields = append(missingFields, "SMTP_USERNAME")
	}
	if a.SMTP.From == "" {
		missingFields = append(missingFields, "SMTP_FROM")
	}
	if a.JWT.Secret == "" {
		missingFields = append(missingFields, "JWT_SECRET")
	}
	if a.OAuth.GoogleClientID == "" {
		missingFields = append(missingFields, "GOOGLE_CLIENT_ID")
	}
	if a.OAuth.GoogleClientSecret == "" {
		missingFields = append(missingFields, "GOOGLE_CLIENT_SECRET")
	}
	if a.OAuth.GithubClientID == "" {
		missingFields = append(missingFields, "GITHUB_CLIENT_ID")
	}
	if a.OAuth.GithubClientSecret == "" {
		missingFields = append(missingFields, "GITHUB_CLIENT_SECRET")
	}
	if a.OAuth.RedirectURL == "" {
		missingFields = append(missingFields, "REDIRECT_URL")
	}

	if len(missingFields) > 0 {
		a.Logger.Error("Missing required fields in configuration", zap.Strings("missingFields", missingFields))
		return fmt.Errorf("missing required fields: %s", missingFields)
	}
	return nil
}
