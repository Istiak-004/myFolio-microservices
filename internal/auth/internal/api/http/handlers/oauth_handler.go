package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/ports"
	oauth_service "github.com/istiak-004/myFolio-microservices/auth/internal/domain/service/oauth"
)

type OAuthHandler struct {
	OauthService ports.OAuthService // supports RegisterOrLoginGoogle(ctx, email, name)
}

func NewOAuthHandler(r *gin.Engine, oauthService ports.OAuthService) {
	h := &OAuthHandler{OauthService: oauthService}

	group := r.Group("/auth/google")
	group.GET("/login", h.Login)
	group.GET("/callback", h.Callback)
}

func (h *OAuthHandler) Login(c *gin.Context) {
	url := oauth_service.GoogleOAuthConfig.AuthCodeURL("random-state") // Add CSRF token later
	c.Redirect(http.StatusFound, url)
}

func (h *OAuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	ctx := c.Request.Context()
	token, err := oauth_service.ExchangeCode(ctx, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "code exchange failed"})
		return
	}

	user, err := oauth_service.GetGoogleUser(ctx, token)
	if err != nil || !user.EmailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or unverified Google account"})
		return
	}

	authToken, err := h.OauthService.RegisterOrLoginGoogle(ctx, user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth service error"})
		return
	}

	setSecureRefreshCookie(c, authToken.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": authToken.AccessToken, "expires_in": authToken.ExpiresIn})
}
