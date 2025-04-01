package handler

import (
	"encoding/json"
	"net/http"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/go-chi/chi/v5"
)

type courseCategoryHandler struct {
	log          logger.Logger
	categoryRepo repository.CourseCategoryRepository
}

func NewCourseCategoryHandler(log logger.Logger, category repository.CourseCategoryRepository) *courseCategoryHandler {
	return &courseCategoryHandler{
		log:          log,
		categoryRepo: category,
	}
}

// GetCourseCategory retrieves a specific course category given the ID.
// @Summary Retrieves a course category.
// @Tags Course Category
// @Param categoryId path string true "Category ID"
// @Failure 500 {string} string
// @Security Session
// @Router /courses/categories/{categoryId} [get]
func (p *courseCategoryHandler) GetCourseCategory(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	categoryId := chi.URLParam(r, "categoryId")
	iCourseCategory, err := p.categoryRepo.GetCourseCategory(r.Context(), categoryId)
	if err != nil {
		http.Error(w, "Failed to get course category", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.ToCourseCategory(iCourseCategory))
}

// ListCourseCategories returns a list of course categories, optionally filtered by search.
// @Summary List course categories
// @Tags Course Category
// @Param search query string false "Search keyword"
// @Success 200 {array} model.CourseCategory
// @Failure 500 {string} string
// @Security Session
// @Router /courses/categories [get]
func (p *courseCategoryHandler) ListCourseCategories(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	query := r.URL.Query()
	iCourseCategories, err := p.categoryRepo.ListCourseCategories(r.Context(), query.Get("search"))
	if err != nil {
		http.Error(w, "Failed to get course categories", http.StatusInternalServerError)
		return
	}

	courseCategories := make([]model.CourseCategory, 0, len(iCourseCategories))
	for _, iCourseCategory := range iCourseCategories {
		courseCategories = append(courseCategories, model.ToCourseCategory(iCourseCategory))
	}

	json.NewEncoder(w).Encode(courseCategories)
}
