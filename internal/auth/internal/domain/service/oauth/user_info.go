package oauth_service

import (
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

type GoogleUser struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
	Sub           string `json:"sub"`
}

func GetGoogleUser(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
