package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash" validate:"required,min=8"`
	GoogleID     string    `json:"-" db:"google_id"`
	GitHubID     string    `json:"-" db:"github_id"`
	Role         string    `json:"role" db:"role" validate:"oneof=admin visitor super_admin"`
	FirstName    string    `json:"first_name,omitempty" db:"first_name"`
	LastName     string    `json:"last_name,omitempty" db:"last_name"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	IsSuperAdmin bool      `json:"is_super_admin" db:"is_super_admin"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
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
