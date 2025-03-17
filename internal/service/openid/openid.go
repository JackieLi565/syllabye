package openid

import "golang.org/x/oauth2"

// https://openid.net/specs/openid-connect-core-1_0.html 5.1
type StandardClaims struct {
	Name          string
	Email         string
	EmailVerified bool
	Picture       string
	Sub           string
}

type OpenIdProvider interface {
	AuthConsentUrl(payload *StateClaims) (string, error)
	VerifyCodeExchange(code string) (*oauth2.Token, error)
	ParseStandardClaims(tokenString string) (StandardClaims, error)
	ParseStateClaims(tokenString string) (*StateClaims, error)
}
