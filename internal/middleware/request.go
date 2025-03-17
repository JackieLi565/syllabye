package middleware

import (
	"context"
	"net/http"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/google/uuid"
)

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New().String()
		w.Header().Set("Request-Id", requestId)

		ctx := context.WithValue(r.Context(), config.RequestIdKey, requestId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
