package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/openid"
)

type AuthHandler struct {
	OpenIdProvider openid.OpenIdProvider
	UserRepo       repository.UserRepository
	SessionRepo    repository.SessionRepository
}

const defaultRedirect = "/home"

func (ah AuthHandler) ConsentUrlHandler(w http.ResponseWriter, r *http.Request) {
	redirectUrl := r.URL.Query().Get("redirect")
	parsedRedirectUrl, err := url.Parse(redirectUrl)
	if err != nil {
		log.Println("invalid redirect url provided")
		redirectUrl = ""
	} else {
		// Redirect must be from the same browser host. In this case its the syllabye domain
		if parsedRedirectUrl.Host != r.Host {
			redirectUrl = ""
		}
	}

	sessionCookie, err := r.Cookie(config.SessionCookie)
	if err == nil {
		session, err := ah.SessionRepo.FindSession(sessionCookie.Value)
		if err == nil {
			if session.DateExpires.After(time.Now()) {
				if redirectUrl == "" {
					http.Redirect(w, r, defaultRedirect, http.StatusFound) // TODO: Change default redirect
				} else {
					http.Redirect(w, r, redirectUrl, http.StatusFound)
				}
				return
			}
		}
	}

	url, err := ah.OpenIdProvider.AuthConsentUrl(&openid.StateClaims{
		Redirect: redirectUrl,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorMessage{
			Message: "Unable to continue to OpenID provider.",
		})
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (ah AuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	state := query.Get("state")

	// Validate state token
	stateClaims, err := ah.OpenIdProvider.ParseStateClaims(state)
	if err != nil {
		log.Println("failed to parse state claim: \n%w", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorMessage{
			Message: "Login state flow no longer valid.",
		})
		return
	}

	// Validate exchange code
	tokens, err := ah.OpenIdProvider.VerifyCodeExchange(code)
	if err != nil {
		log.Println("failed to exchange code for tokens: \n%w", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorMessage{
			Message: "An error occurred while exchanging for authorization tokens.",
		})
		return
	}

	// Validate extra token
	idToken, ok := tokens.Extra("id_token").(string)
	if !ok {
		log.Println("failed to type assert id token: \n%w", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorMessage{
			Message: "Unable to validate ID token.",
		})
		return
	}

	// Validate OpenID token and email domain
	standardClaims, err := ah.OpenIdProvider.ParseStandardClaims(idToken)
	splitEmail := strings.Split(standardClaims.Email, "@")
	if len(splitEmail) != 2 {
		log.Printf("unknown email format received from open id %s\n", standardClaims.Email)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorMessage{
			Message: "Unable to validate email.",
		})
		return
	}
	if splitEmail[1] != "torontomu.ca" {
		log.Printf("unauthorized email login attempt %s\n", standardClaims.Email)
		http.Redirect(w, r, "/sorry", http.StatusFound) // TODO: unauthorized route
		return
	}

	// Login or Register the user
	user, err := ah.UserRepo.LoginOrRegisterUser(standardClaims)
	if err != nil {
		log.Println("failed to login or register user: \n%w", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorMessage{
			Message: "Unable to login or register user.",
		})
		return
	}

	// TODO: current implementation of a session token lasts 30 days. Implement refresh token in the future
	sessionExp := time.Now().Add(720 * time.Hour)
	session, err := ah.SessionRepo.CreateSession(user.Id, sessionExp)
	if err != nil {
		log.Println("failed to create session for user: \n%w", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorMessage{
			Message: "Unable to create session for user.",
		})
		return
	}
	log.Println(stateClaims.Redirect)
	http.SetCookie(w, &http.Cookie{
		Name:     config.SessionCookie,
		Value:    session.Id,
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

func (ah AuthHandler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(config.SessionCookie)
		if err == nil {
			session, err := ah.SessionRepo.FindSession(sessionCookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(model.ErrorMessage{
					Message: "Session not found.",
				})
				return
			}

			ctx := context.WithValue(r.Context(), config.SessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
