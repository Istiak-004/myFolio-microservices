package models

import "time"

type OAuthUser struct {
	ID    string `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	Name  string `json:"name" db:"name"`
}

type OauthProviders struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Provider   string    `json:"provider"`
	ProviderID string    `json:"provider_id"`
	CreatedAt  time.Time `json:"created_at"`
}
