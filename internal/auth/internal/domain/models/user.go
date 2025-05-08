package models

import "time"

type User struct {
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash" validate:"required,min=8"`
	GoogleID     string    `json:"-"`
	GitHubID     string    `json:"-"`
	Role         string    `json:"role" validate:"oneof=admin visitor super_admin" db:"role"`
	FirstName    string    `json:"first_name,omitempty" db:"first_name"`
	LastName     string    `json:"last_name,omitempty" db:"last_name"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	IsSuperAdmin bool      `json:"is_super_admin" db:"is_super_admin"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type VerificationToken struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

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

// EmailVerification represents an email verification request
type EmailVerification struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	UsedAt    time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
