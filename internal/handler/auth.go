package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/authorizer"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/service/openid"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

const defaultRedirectUri = "/"

type authHandler struct {
	log            logger.Logger
	openIdProvider openid.OpenIdProvider
	userRepo       repository.UserRepository
	sessionRepo    repository.SessionRepository
	jwt            *authorizer.JwtAuthorizer
}

func NewAuthHandler(log logger.Logger, user repository.UserRepository, session repository.SessionRepository, openId openid.OpenIdProvider, jwt *authorizer.JwtAuthorizer) *authHandler {
	return &authHandler{
		log:            log,
		openIdProvider: openId,
		userRepo:       user,
		sessionRepo:    session,
		jwt:            jwt,
	}
}

// ConsentUrlRedirect initiates the login flow by redirecting to the OpenID provider's consent screen.
// @Summary Redirect to OpenID consent screen
// @Description Validates an optional redirect query param and redirects the user to the OpenID login flow.
// @Tags Authentication
// @Param redirect query string false "Optional redirect URL after login"
// @Success 302 {string} string "Redirects to OpenID consent screen"
// @Failure 500 {string} string "Unable to continue to OpenID provider"
// @Router /providers/google [get]
func (ah *authHandler) ConsentUrlRedirect(w http.ResponseWriter, r *http.Request) {
	redirectUrl := ah.getValidRedirectUrl(r.URL.Query().Get("redirect"), os.Getenv(config.Domain))
	ah.log.Info(redirectUrl)

	sessionCookie, err := r.Cookie(config.SessionCookie)
	if err == nil {
		_, err := ah.decodeSessionToken(sessionCookie.Value)
		if err == nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return
		}
	}

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

// ProviderCallback handles the OAuth2 internal callback from the OpenID provider.
// ProviderCallback handles the OAuth2 internal callback from the OpenID provider.
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
		http.Redirect(w, r, os.Getenv(config.ClientDomain)+"/sorry", http.StatusFound)
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

	sessionToken, err := ah.jwt.EncodeJwt(jwt.MapClaims{
		"id":     sessionId,
		"userId": userId,
	})

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
		http.Redirect(w, r, os.Getenv(config.ClientDomain)+defaultRedirectUri, http.StatusFound)
	}

	ah.log.Info(fmt.Sprintf("user %s has logged in with session token %s", userId, sessionToken))
}

// Logout removes/invalidates the user's session cookie
// @Summary Logout user session
// @Description Removes the users session cookie if exists.
// @Tags Authentication
// @Success 302 {string} string "Redirects to root page"
// @Router /logout [get]
func (ah *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.SessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, os.Getenv(config.ClientDomain)+defaultRedirectUri, http.StatusFound)
	ah.log.Info("a user has logged out")
}

// SessionCheck verifies the user's session cookie and returns session information if valid.
// @Summary Check user session
// @Description Validates the session cookie and returns session payload if authenticated.
// @Tags Authentication
// @Success 200 {object} model.Session "Valid session"
// @Failure 401 {string} string "Missing or invalid session cookie"
// @Router /me [get]
// @Security Session
func (ah *authHandler) SessionCheck(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

// AuthMiddleware secures endpoints with both session and bearer token authorization.
func (ah *authHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for session auth first
		sessionCookie, err := r.Cookie(config.SessionCookie)
		if err != nil {
			// Cookie not found, check for authorization bearer token
			token, err := util.GetBearerToken(r)
			if err != nil {
				http.Error(w, "Session or authorization token not found.", http.StatusUnauthorized)
				return
			}

			if _, err = ah.jwt.DecodeJwt(token); err != nil {
				ah.log.Info("invalid authorization token")
				http.Error(w, "Internal route access denied.", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		session, err := ah.decodeSessionToken(sessionCookie.Value)
		if err != nil {
			http.Error(w, "Invalid session token.", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), config.AuthKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// decodeSessionToken decodes a token string to a session model.
func (ah *authHandler) decodeSessionToken(tokenString string) (model.Session, error) {
	claims, err := ah.jwt.DecodeJwt(tokenString)
	if err != nil {
		ah.log.Info("failed to decode jwt", logger.Err(err))
		return model.Session{}, err
	}

	var session model.Session
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

	return session, nil
}

// getValidRedirectUrl retrieves a valid redirect url, otherwise a default url is provided.
func (ah *authHandler) getValidRedirectUrl(redirectUrl string, host string) string {
	defaultUrl := os.Getenv(config.ClientDomain) + defaultRedirectUri
	if redirectUrl == "" {
		return defaultUrl
	}

	parsedRedirectUrl, err := url.Parse(redirectUrl)
	if err != nil {
		ah.log.Info("invalid redirect url provided")
		return defaultUrl
	}

	if parsedRedirectUrl.Host != host {
		ah.log.Info(fmt.Sprintf("restricted redirect url %s", parsedRedirectUrl.String()))
		return defaultUrl
	}

	return parsedRedirectUrl.String()
}

func (ah *authHandler) DevAuthorization(w http.ResponseWriter, r *http.Request) {
	if os.Getenv(config.ENV) != "development" {
		ah.log.Error("development authorization handler used outside of development")
		http.Error(w, "We are unable to handle your login request at this time. Please try again next time.", http.StatusInternalServerError)
		return
	}

	redirectUrl := ah.getValidRedirectUrl(r.URL.Query().Get("redirect"), r.Host)

	userId, err := ah.userRepo.LoginOrRegisterUser(r.Context(), openid.StandardClaims{
		Name:  "System Admin (Development)",
		Email: "sys.admin@syllabye.ca",
	})
	if err != nil {
		http.Error(w, "Dev authorization failed.", http.StatusInternalServerError)
		return
	}

	sessionToken, err := ah.jwt.EncodeJwt(jwt.MapClaims{
		"id":     "", // Only allowed in dev handler
		"userId": userId,
	})
	if err != nil {
		ah.log.Error("failed to encode session token", logger.Err(err))
		http.Error(w, "Unable to create admin session.", http.StatusInternalServerError)
		return
	}

	sessionExp := time.Now().Add(time.Hour * 24 * 30 * 12) // 1 Year
	cookie := &http.Cookie{
		Name:     config.SessionCookie,
		Value:    sessionToken,
		Path:     "/",
		Expires:  sessionExp,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}
