package http

import (
	"time"

	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/models"
)

// RegisterRequest defines the request body for registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest defines the request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse defines the token response
type TokenResponse struct {
	ExpiresAt    time.Time `json:"expires_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
}

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	u:=models.User{}
	println(u)
	return &AuthHandler{}
}
