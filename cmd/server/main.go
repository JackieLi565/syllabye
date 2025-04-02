package main

import (
	"net/http"
	"os"

	_ "github.com/JackieLi565/syllabye/docs"
	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/handler"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/bucket"
	"github.com/JackieLi565/syllabye/internal/service/database"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/service/openid"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Syllabye API
// @version 1.0
// @description Syllabye API server.

// @contact.name Jackie Li

// @BasePath /api

// @securityDefinitions.apiKey Session
// @in cookie
// @name syllabye.session

// @securityDefinitions.apikey AWS
// @in header
// @name Authorization
// @description Provide a valid Bearer token. Example: "Bearer {jwt-token}"
func main() {
	env := os.Getenv(config.ENV)

	var log logger.Logger
	if env == "production" {
		log = logger.NewJsonLogger()
	} else {
		log = logger.NewTextLogger()
	}

	// Services
	db, err := database.NewPostgresDb()
	if err != nil { // TODO: remove panic and panic from function level
		panic("database connection failed")
	}
	s3Client := bucket.NewS3Client(log)

	googleOpenId := openid.NewGoogleOpenIdProvider(log)
	s3Presigner := bucket.NewS3Presigner(log, s3Client, os.Getenv(config.AWS_S3_SYLLABI_BUCKET))

	// Repositories
	pgProgramRepo := repository.NewPgProgramRepository(db, log)
	pgSessionRepo := repository.NewPgSessionRepository(db, log)
	pgUserRepo := repository.NewPgUserRepository(db, log)
	pgFacultyRepo := repository.NewPgFacultyRepository(db, log)
	pgCourseCategoryRepo := repository.NewPgCourseCategoryRepository(db, log)
	pgCourseRepo := repository.NewPgCourseRepository(db, log)
	pgSyllabusRepo := repository.NewPgSyllabusRepository(db, log)

	// Handlers
	utilHandler := handler.NewUtilHandler()
	authHandler := handler.NewAuthHandler(log, pgUserRepo, pgSessionRepo, googleOpenId)
	programHandler := handler.NewProgramHandler(log, pgProgramRepo)
	facultyHandler := handler.NewFacultyHandler(log, pgFacultyRepo)
	courseCategoryHandler := handler.NewCourseCategoryHandler(log, pgCourseCategoryRepo)
	courseHandler := handler.NewCourseHandler(log, pgCourseRepo)
	userHandler := handler.NewUserHandler(log, pgUserRepo)
	syllabusHandler := handler.NewSyllabusHandler(log, pgSyllabusRepo, s3Presigner)

	r := chi.NewRouter()
	r.Use(utilHandler.RequestIdMiddleware)

	var basePath string
	if env == "development" {
		basePath = "/api"

		r.Use(utilHandler.AllowAllCORSMiddleware)

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
	} else {
		basePath = "/"
	}

	r.Route(basePath, func(r chi.Router) {
		r.Route("/providers/google", func(r chi.Router) {
			r.Get("/", authHandler.ConsentUrlRedirect)
			r.Get("/callback", authHandler.ProviderCallback)
		})

		r.Route("/me", func(r chi.Router) {
			r.Use(utilHandler.JsonMiddleware)

			r.Get("/", authHandler.SessionCheck)
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

			r.Route("/{userId}", func(r chi.Router) {
				r.Get("/", userHandler.GetUser)
				r.Patch("/", userHandler.UpdateUser)
			})
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

		r.Route("/syllabi", func(r chi.Router) {
			r.Use(authHandler.SessionMiddleware)
			r.Use(utilHandler.JsonMiddleware)

			r.Post("/", syllabusHandler.CreateSyllabus)
			r.Get("/", syllabusHandler.ListSyllabi)

			r.Route("/{syllabusId}", func(r chi.Router) {
				r.Get("/", syllabusHandler.GetSyllabus)
				r.Patch("/", syllabusHandler.UpdateSyllabus)
				r.Delete("/", syllabusHandler.DeleteSyllabus)
				r.Get("/sync", syllabusHandler.SyncSyllabus)

				r.Route("/reactions", func(r chi.Router) {
					r.Get("/", syllabusHandler.ListSyllabusLikes)

					r.Post("/", syllabusHandler.SyllabusReaction)
					r.Delete("/", syllabusHandler.DeleteSyllabusReaction)
				})
			})
		})
	})

	http.ListenAndServe(os.Getenv("PORT"), r)
}
