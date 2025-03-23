package handler

import (
	"context"
	"net/http"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/google/uuid"
)

type utilHandler struct{}

func NewUtilHandler() *utilHandler {
	return &utilHandler{}
}

func (u *utilHandler) JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (U *utilHandler) RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New().String()
		w.Header().Set("Request-Id", requestId)

		ctx := context.WithValue(r.Context(), config.RequestIdKey, requestId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
