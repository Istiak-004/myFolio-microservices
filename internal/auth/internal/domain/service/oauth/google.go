package oauth_service

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig = &oauth2.Config{
	ClientID:     "", // <- set via env later
	ClientSecret: "",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return GoogleOAuthConfig.Exchange(ctx, code)
}
