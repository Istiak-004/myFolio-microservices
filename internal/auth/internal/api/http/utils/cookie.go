package utils

import (
	"net/http"
)

const (
	RefreshTokenCookieName = "refresh_token"
	RefreshTokenPath       = "/auth/refresh"  // Only send on this path
	RefreshTokenMaxAge     = 7 * 24 * 60 * 60 // 7 days
)

// SetRefreshTokenCookie sets the refresh token securely as an HTTP-only cookie
func SetRefreshTokenCookie(w http.ResponseWriter, token string, domain string) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    token,
		Path:     RefreshTokenPath,
		Domain:   domain, // e.g. "myapp.com"
		MaxAge:   RefreshTokenMaxAge,
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearRefreshTokenCookie deletes the cookie
func ClearRefreshTokenCookie(w http.ResponseWriter, domain string) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     RefreshTokenPath,
		Domain:   domain,
		MaxAge:   -1, // Delete immediately
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
	})
}
