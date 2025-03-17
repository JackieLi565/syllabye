package openid

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type StateClaims struct {
	jwt.RegisteredClaims
	Redirect string                 `json:"redirect,omitempty"`
	State    map[string]interface{} `json:"state,omitempty"`
}

func parseStateClaims(tokenString string) (*StateClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &StateClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv(config.JwtSecret)), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token not valid")
	}

	if claims, ok := token.Claims.(*StateClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("failed to type assert claims")
	}
}

func newStateClaims(payload *StateClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"redirect": payload.Redirect,
		"state":    payload.State,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
		"iss":      config.JwtIssuer,
		"aud":      config.JwtIssuer,
	})

	jwtSecret := os.Getenv(config.JwtSecret)
	if jwtSecret == "" {
		log.Fatal("Invalid JWT signing secret")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
