package internal

import (
	"github.com/JackieLi565/syllabye/internal/handler"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Auth handler.AuthHandler
}

func (ro Router) SetupRoutes(r chi.Router) {
	r.Route("/providers/google", func(r chi.Router) {
		// All redirect routes
		r.Get("/", ro.Auth.ConsentUrlHandler)
		r.Get("/callback", ro.Auth.CallbackHandler)
	})
}
