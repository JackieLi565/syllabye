package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/service/openid"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

const defaultRedirect = "/home"

type authHandler struct {
	log            logger.Logger
	openIdProvider openid.OpenIdProvider
	userRepo       repository.UserRepository
	sessionRepo    repository.SessionRepository
}

func NewAuthHandler(log logger.Logger, user repository.UserRepository, session repository.SessionRepository, openId openid.OpenIdProvider) *authHandler {
	return &authHandler{
		log:            log,
		openIdProvider: openId,
		userRepo:       user,
		sessionRepo:    session,
	}
}

func (ah *authHandler) ConsentUrlRedirect(w http.ResponseWriter, r *http.Request) {
	redirectUrl := r.URL.Query().Get("redirect")
	parsedRedirectUrl, err := url.Parse(redirectUrl)
	if err != nil {
		ah.log.Warn("invalid redirect url provided")
		redirectUrl = ""
	} else {
		// Redirect must be from the same browser host. In this case its the syllabye domain
		if parsedRedirectUrl.Host != r.Host {
			redirectUrl = ""
		}
	}

	sessionCookie, err := r.Cookie(config.SessionCookie)
	if err == nil {
		_, err := ah.decodeSessionToken(sessionCookie.Value)
		if err == nil {
			http.Redirect(w, r, defaultRedirect, http.StatusFound)
			return
		}
	}
	ah.log.Info(redirectUrl)
	url, err := ah.openIdProvider.AuthConsentUrl(&openid.StateClaims{
		Redirect: redirectUrl,
	})
	if err != nil {
		ah.log.Warn("failed to continue to login provider", logger.Err(err))
		http.Error(w, "Unable to continue to OpenID provider.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (ah *authHandler) ProviderCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	state := query.Get("state")

	// Validate state token
	stateClaims, err := ah.openIdProvider.ParseStateClaims(state)
	if err != nil {
		ah.log.Warn("login state claim expired", logger.Err(err))
		http.Error(w, "Login state flow no longer valid.", http.StatusUnauthorized)
		return
	}

	// Get OpenID token from exchange code
	idToken, err := ah.openIdProvider.GetOpenIdToken(code)
	if err != nil {
		// TODO: improve error handling response code. Not all errors result in a internal error
		http.Error(w, "Unable to validate ID token.", http.StatusInternalServerError)
		return
	}

	// Validate OpenID token and email domain
	standardClaims, err := ah.openIdProvider.ParseStandardClaims(idToken)
	splitEmail := strings.Split(standardClaims.Email, "@")
	if len(splitEmail) != 2 {
		ah.log.Error("unknown email format received from open id")
		http.Error(w, "Unable to validate ID token.", http.StatusInternalServerError)
		return
	}
	if splitEmail[1] != "torontomu.ca" {
		ah.log.Info(fmt.Sprintf("unauthorized login attempt with email %s", standardClaims.Email))
		http.Redirect(w, r, "/sorry", http.StatusFound) // TODO: unauthorized route
		return
	}

	// Login or Register the user
	userId, err := ah.userRepo.LoginOrRegisterUser(r.Context(), standardClaims)
	if err != nil {
		ah.log.Error("failed to login or register user", logger.Err(err))
		http.Error(w, "Unable to login or register user.", http.StatusInternalServerError)
		return
	}

	sessionId, err := ah.sessionRepo.CreateSession(r.Context(), userId) // create session log in database
	if err != nil {
		ah.log.Error("failed to create session log", logger.Err(err))
		http.Error(w, "Unable to create session for user.", http.StatusInternalServerError)
		return
	}

	sessionToken, err := ah.encodeSessionToken(userId, sessionId)
	if err != nil {
		ah.log.Error("failed to encode session token", logger.Err(err))
		http.Error(w, "Unable to create session for user.", http.StatusInternalServerError)
		return
	}

	sessionExp := time.Now().Add(time.Hour * 24 * 30)
	http.SetCookie(w, &http.Cookie{
		Name:     config.SessionCookie,
		Value:    sessionToken,
		Path:     "/",
		Expires:  sessionExp,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	if stateClaims.Redirect != "" {
		http.Redirect(w, r, stateClaims.Redirect, http.StatusFound)
	} else {
		http.Redirect(w, r, defaultRedirect, http.StatusFound)
	}
}

func (ah *authHandler) SessionCheck(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		ah.log.Error("missing session middleware")
	}

}

func (ah *authHandler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(config.SessionCookie)
		if err != nil {
			http.Error(w, "Cookie not found.", http.StatusUnauthorized)
			return
		}

		session, err := ah.decodeSessionToken(sessionCookie.Value)
		if err != nil {
			http.Error(w, "Invalid session token.", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), config.SessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ah *authHandler) encodeSessionToken(userId string, sessionId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     sessionId,
		"userId": userId,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv(config.JwtSecret)))
	if err != nil {
		ah.log.Info("failed to sign jwt session token")
		return "", err
	}

	return tokenString, nil
}

func (ah *authHandler) decodeSessionToken(tokenString string) (model.Session, error) {
	var session model.Session

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			ah.log.Info(fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]))
			return nil, util.ErrMalformed
		}

		return []byte(os.Getenv(config.JwtSecret)), nil
	})
	if err != nil {
		return model.Session{}, err
	}

	if !token.Valid {
		ah.log.Info("invalid jwt token")
		return model.Session{}, util.ErrMalformed
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims["userId"].(string); ok {
			session.UserId = userId
		} else {
			ah.log.Info("userId not found within parsed claim")
			return model.Session{}, util.ErrMalformed
		}

		if id, ok := claims["id"].(string); ok {
			session.Id = id
		} else {
			ah.log.Info("id not found within parsed claim")
			return model.Session{}, util.ErrMalformed
		}
	}

	return session, nil
}
