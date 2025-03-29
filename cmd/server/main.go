package main

import (
	"net/http"
	"os"

	_ "github.com/JackieLi565/syllabye/docs"
	"github.com/JackieLi565/syllabye/internal/handler"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/service/openid"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

// @title Syllabye API
// @version 1.0
// @description Syllabye API server.

// @contact.name Jackie Li

// @securityDefinitions.apiKey Session
// @in cookie
// @name syllabye.session
func main() {
	env := os.Getenv("ENV")

	var log logger.Logger
	if env == "production" {
		log = logger.NewJsonLogger()
	} else {
		log = logger.NewTextLogger()
	}

	// Globals
	db, err := database.NewPostgresDb()
	if err != nil {
		panic("database connection failed")
	}

	// Repositories
	pgProgramRepo := repository.NewPgProgramRepository(db, log)
	pgSessionRepo := repository.NewPgSessionRepository(db, log)
	pgUserRepo := repository.NewPgUserRepository(db, log)
	pgFacultyRepo := repository.NewPgFacultyRepository(db, log)
	pgCourseCategoryRepo := repository.NewPgCourseCategoryRepository(db, log)
	pgCourseRepo := repository.NewPgCourseRepository(db, log)

	// Services
	googleOpenId := openid.NewGoogleOpenIdProvider(log)

	// Handlers
	utilHandler := handler.NewUtilHandler()
	authHandler := handler.NewAuthHandler(log, pgUserRepo, pgSessionRepo, googleOpenId)
	programHandler := handler.NewProgramHandler(log, pgProgramRepo)
	facultyHandler := handler.NewFacultyHandler(log, pgFacultyRepo)
	courseCategoryHandler := handler.NewCourseCategoryHandler(log, pgCourseCategoryRepo)
	courseHandler := handler.NewCourseHandler(log, pgCourseRepo)
	userHandler := handler.NewUserHandler(log, pgUserRepo)

	r := chi.NewRouter()
	r.Use(utilHandler.RequestIdMiddleware)

	if env == "development" {
		r.Route("/openapi", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/openapi/index.html", http.StatusFound)
			})

			r.Get("/doc.json", func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, "./docs/swagger.json")
			})

			r.Get("/*", httpSwagger.Handler(
				httpSwagger.URL("http://localhost:8000/openapi/doc.json"),
			))
		})

		r.Route("/api", func(r chi.Router) {
			r.Route("/providers/google", func(r chi.Router) {
				r.Get("/", authHandler.ConsentUrlRedirect)
				r.Get("/callback", authHandler.ProviderCallback)
			})

			r.Route("/programs", func(r chi.Router) {
				r.Use(authHandler.SessionMiddleware)
				r.Use(utilHandler.JsonMiddleware)

				r.Get("/", programHandler.ListPrograms)
				r.Get("/{programId}", programHandler.GetProgram)
			})

			r.Route("/faculties", func(r chi.Router) {
				r.Use(authHandler.SessionMiddleware)
				r.Use(utilHandler.JsonMiddleware)

				r.Get("/", facultyHandler.ListFaculties)
				r.Get("/{facultyId}", facultyHandler.GetFaculty)
			})

			r.Route("/users", func(r chi.Router) {
				r.Use(authHandler.SessionMiddleware)
				r.Use(utilHandler.JsonMiddleware)

				r.Get("/{userId}", userHandler.GetUser)
				r.Patch("/{userId}", userHandler.UpdateUser)
			})

			r.Route("/courses", func(r chi.Router) {
				r.Use(authHandler.SessionMiddleware)
				r.Use(utilHandler.JsonMiddleware)

				r.Get("/", courseHandler.ListCourses)
				r.Get("/{courseId}", courseHandler.GetCourse)

				r.Route("/categories", func(r chi.Router) {
					r.Get("/", courseCategoryHandler.ListCourseCategories)
					r.Get("/{categoryId}", courseCategoryHandler.GetCourseCategory)
				})
			})
		})
	}

	http.ListenAndServe(os.Getenv("PORT"), r)
}
