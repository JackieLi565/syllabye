package openid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/JackieLi565/syllabye/server/internal/config"
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
type GoogleOpenIdProvider struct{}

func NewGoogleOpenIdProvider() *GoogleOpenIdProvider {
	return &GoogleOpenIdProvider{}
}

func (g GoogleOpenIdProvider) NewConsentUrl(payload *StateClaims) (string, error) {
	stateToken, err := newStateClaims(payload)
	if err != nil {
		return "", err
	}

	return config.GoogleOAuthConfig.AuthCodeURL(stateToken), nil
}

func (g GoogleOpenIdProvider) VerifyTokenExchange(code string) (*TokenExchangeResponse, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", config.GoogleOAuthConfig.ClientID)
	form.Set("client_secret", config.GoogleOAuthConfig.ClientSecret)
	form.Set("redirect_uri", config.GoogleOAuthConfig.RedirectURL)
	form.Set("grant_type", "authorization_code")

	resp, err := http.Post(config.GoogleOAuthConfig.Endpoint.TokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("authorization server post request failed")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse TokenExchangeResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (g GoogleOpenIdProvider) ParseStandardClaims(tokenString string) (*StandardClaims, error) {
	token, err := idtoken.Validate(context.Background(), tokenString, config.GoogleOAuthConfig.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Google ID token: %w", err)
	}

	claimsJson, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize Google ID claims")
	}

	var claims googleStandardClaims
	err = json.Unmarshal(claimsJson, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Google ID claims")
	}

	return &StandardClaims{
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
