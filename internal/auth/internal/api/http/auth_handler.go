package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/ports"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/valueobjects"
	"github.com/istiak-004/myFolio-microservices/pkg/http/middleware"
)

// TokenResponse defines the token response
type TokenResponse struct {
	ExpiresAt    time.Time `json:"expires_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
}

const (
	RefreshTokenCookieName = "refresh_token"
)

type AuthHandler struct {
	AuthService ports.AuthService
}

func NewAuthHandler(r *gin.Engine, authService ports.AuthService) {
	h := &AuthHandler{authService}

	group := r.Group("/auth")
	group.Use(middleware.RateLimiterMiddleware(middleware.RateLimitConfig{
		Enabled:  true,
		Requests: 10,
		Interval: time.Minute,
	}), middleware.SecurityHeaders())

	group.POST("/register", h.HandleRegister)
	group.POST("/login", h.HandleLogin)
	group.POST("/refresh", h.HandleRefresh)
	group.POST("/logout", h.HandleLogout)
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

func (h *AuthHandler) HandleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email, err := valueobjects.NewEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	password, err := valueobjects.NewPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// register the user
	_, err = h.AuthService.Register(c.Request.Context(), email, password, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registration successful, please check your email"})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) HandleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email, err := valueobjects.NewEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	password, err := valueobjects.NewPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := h.AuthService.Login(c.Request.Context(), email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	setSecureRefreshCookie(c, tokens.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "expires_in": tokens.ExpiresIn})
}

func (h *AuthHandler) HandleRefresh(c *gin.Context) {
	cookie, err := c.Cookie(RefreshTokenCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}
	token := valueobjects.Token{TokenString: cookie}
	tokens, err := h.AuthService.RefreshToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}
	setSecureRefreshCookie(c, tokens.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "expires_in": tokens.ExpiresIn})
}

func (h *AuthHandler) HandleLogout(c *gin.Context) {
	cookie, err := c.Cookie(RefreshTokenCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}
	token := valueobjects.Token{TokenString: cookie}
	if err := h.AuthService.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}
	clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func setSecureRefreshCookie(c *gin.Context, refreshToken string) {
	c.SetCookie(
		RefreshTokenCookieName,
		refreshToken,
		7*24*3600, // 7 days
		"/",
		"",   // domain
		true, // Secure
		true, // HttpOnly
	)
	c.Writer.Header().Add("Set-Cookie", "SameSite=Strict")
}

func clearRefreshCookie(c *gin.Context) {
	c.SetCookie(
		RefreshTokenCookieName,
		"",
		-1,
		"/",
		"",
		true,
		true,
	)
	c.Writer.Header().Add("Set-Cookie", "SameSite=Strict")
}
