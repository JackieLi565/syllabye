package main

import (
	"log"
	"net/http"
	"os"

	"github.com/JackieLi565/syllabye/internal"
	"github.com/JackieLi565/syllabye/internal/handler"
	"github.com/JackieLi565/syllabye/internal/middleware"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/openid"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	pgUserRepo := repository.PgUserRepository{
		DB: db,
	}
	pgSessionRepo := repository.PgSessionRepository{
		DB: db,
	}

	googleOpenId := openid.NewGoogleOpenIdProvider()

	authHandler := handler.AuthHandler{
		OpenIdProvider: googleOpenId,
		UserRepo:       pgUserRepo,
		SessionRepo:    pgSessionRepo,
	}

	router := internal.Router{
		Auth: authHandler,
	}

	r := chi.NewRouter()
	log.Println(os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"))
	r.Use(middleware.RequestIdMiddleware)

	env := os.Getenv("ENV") // development | production
	if env == "development" {
		r.Route("/api", router.SetupRoutes)
	} else {
		router.SetupRoutes(r)
	}

	log.Println("Sever now listening on http://localhost:8000")
	http.ListenAndServe(":8000", r)

}
