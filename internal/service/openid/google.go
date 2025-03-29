package openid

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

type GoogleOpenIdProvider struct {
	config *oauth2.Config
	log    logger.Logger
}

func NewGoogleOpenIdProvider(log logger.Logger) *GoogleOpenIdProvider {
	return &GoogleOpenIdProvider{
		config: &oauth2.Config{
			ClientID:     os.Getenv(config.GoogleOAuthClientId),
			ClientSecret: os.Getenv(config.GoogleOAuthClientSecret),
			RedirectURL:  os.Getenv(config.GoogleOAuthRedirectUrl),
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"openid",
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
		log: log,
	}
}

func (g *GoogleOpenIdProvider) AuthConsentUrl(payload *StateClaims) (string, error) {
	stateToken, err := newStateClaims(payload)
	if err != nil {
		return "", err
	}

	return g.config.AuthCodeURL(stateToken, oauth2.SetAuthURLParam("hd", "torontomu.ca")), nil
}

func (g *GoogleOpenIdProvider) GetOpenIdToken(code string) (string, error) {
	form := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"client_id":     {g.config.ClientID},
		"client_secret": {g.config.ClientSecret},
		"redirect_uri":  {g.config.RedirectURL},
	}

	resp, err := http.PostForm(g.config.Endpoint.TokenURL, form)
	if err != nil {
		g.log.Error("failed to reach code exchange endpoint", logger.Err(err))
		return "", util.ErrInternal
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		g.log.Error("unsuccessful code exchange request", logger.Err(err))
		return "", util.ErrInternal
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		g.log.Error("failed to parse code exchange response body", logger.Err(err))
		return "", util.ErrMalformed
	}

	// Get id_token
	if idToken, ok := parsed["id_token"].(string); ok {
		return idToken, nil
	} else {
		g.log.Error("code exchange response body did not contain 'id_token' key")
		return "", util.ErrMalformed
	}
}

// VerifyCodeExchange current implementation with oauth2 does not work
// Error invalid_grant Bad Request is thrown with correct credentials
// Most likely internal lib error.
func (g *GoogleOpenIdProvider) VerifyCodeExchange(code string) (*oauth2.Token, error) {
	token, err := g.config.Exchange(context.TODO(), code)
	if err != nil {
		log.Println("code authorization failed")
		return nil, err
	}

	return token, nil
}

func (g *GoogleOpenIdProvider) ParseStandardClaims(tokenString string) (StandardClaims, error) {
	var standardClaims StandardClaims
	token, err := idtoken.Validate(context.TODO(), tokenString, g.config.ClientID)
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

func (g *GoogleOpenIdProvider) ParseStateClaims(tokenString string) (*StateClaims, error) {
	return parseStateClaims(tokenString)
}
