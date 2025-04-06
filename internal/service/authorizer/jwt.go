package authorizer

import (
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

type JwtAuthorizer struct {
	secret string
}

func NewJwtAuthorizer(secret string) *JwtAuthorizer {
	return &JwtAuthorizer{
		secret: secret,
	}
}

// EncodeJwt takes a mapped claim into a jwt token string
func (j *JwtAuthorizer) EncodeJwt(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// DecodeJwt decodes a jwt token string to its mapped claims
func (j *JwtAuthorizer) DecodeJwt(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, util.ErrMalformed
		}

		return []byte(j.secret), nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	if !token.Valid {
		return jwt.MapClaims{}, util.ErrMalformed
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return jwt.MapClaims{}, util.ErrConflict
	}
}
