package openid

// https://openid.net/specs/openid-connect-core-1_0.html 3.1.3.3
type TokenExchangeResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
}

// https://openid.net/specs/openid-connect-core-1_0.html 5.1
type StandardClaims struct {
	Name          string
	Email         string
	EmailVerified bool
	Picture       string
	Sub           string
}

type OpenIdProvider interface {
	NewConsentUrl(payload *StateClaims) (string, error)
	VerifyTokenExchange(code string) (*TokenExchangeResponse, error)
	ParseStandardClaims(tokenString string) (*StandardClaims, error)
	ParseStateClaims(tokenString string) (*StateClaims, error)
}
