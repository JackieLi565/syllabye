package config

import (
	"os"

	"golang.org/x/oauth2"
)

var GoogleOAuthConfig = oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://www.googleapis.com/oauth2/v4/token",
	},
	Scopes: []string{
		"openid",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
}
