package openid

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

// https://developers.google.com/identity/openid-connect/openid-connect

type googleStandardClaims struct {
	HD            string `json:"hd,omitempty"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Sub           string `json:"sub"`
}

// TODO: add logger
type GoogleOpenIdProvider struct {
	config *oauth2.Config
}

func NewGoogleOpenIdProvider() *GoogleOpenIdProvider {
	return &GoogleOpenIdProvider{
		config: &oauth2.Config{
			ClientID:     os.Getenv(config.GoogleOAuthClientId),
			ClientSecret: os.Getenv(config.GoogleOAuthClientSecret),
			RedirectURL:  os.Getenv(config.GoogleOAuthRedirectUrl),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://www.googleapis.com/oauth2/v4/token",
			},
			Scopes: []string{
				"openid",
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
	}
}

func (g GoogleOpenIdProvider) AuthConsentUrl(payload *StateClaims) (string, error) {
	stateToken, err := newStateClaims(payload)
	if err != nil {
		return "", err
	}

	return g.config.AuthCodeURL(stateToken, oauth2.SetAuthURLParam("hd", "torontomu.ca")), nil
}

func (g GoogleOpenIdProvider) VerifyCodeExchange(code string) (*oauth2.Token, error) {
	token, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		log.Println("code authorization failed")
		return nil, err
	}

	return token, nil
}

func (g GoogleOpenIdProvider) ParseStandardClaims(tokenString string) (StandardClaims, error) {
	var standardClaims StandardClaims
	token, err := idtoken.Validate(context.Background(), tokenString, g.config.ClientID)
	if err != nil {
		return standardClaims, fmt.Errorf("failed to verify Google ID token: %w", err)
	}

	claimsJson, err := json.Marshal(token.Claims)
	if err != nil {
		return standardClaims, fmt.Errorf("failed to serialize Google ID claims")
	}

	var claims googleStandardClaims
	err = json.Unmarshal(claimsJson, &claims)
	if err != nil {
		return standardClaims, fmt.Errorf("failed to parse Google ID claims")
	}

	return StandardClaims{
		Name:          claims.Name,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		Picture:       claims.Picture,
		Sub:           claims.Sub,
	}, nil
}

func (g GoogleOpenIdProvider) ParseStateClaims(tokenString string) (*StateClaims, error) {
	return parseStateClaims(tokenString)
}
