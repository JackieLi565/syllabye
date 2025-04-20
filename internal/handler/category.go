package handler

import (
	"encoding/json"
	"net/http"

	"github.com/JackieLi565/syllabye/internal/config"
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

type CourseCategoryRes struct {
	Id   string `json:"id"`
	Name string `json:"name"`
} //@name CourseCategoryResponse

// GetCourseCategory retrieves a specific course category given the ID.
// @Summary Retrieves a course category.
// @Tags Course Category
// @Param categoryId path string true "Category ID"
// @Success 200 {object} CourseCategoryResponse
// @Failure 500 {string} string
// @Security Session
// @Router /courses/categories/{categoryId} [get]
func (p *courseCategoryHandler) GetCourseCategory(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.AuthKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	categoryId := chi.URLParam(r, "categoryId")
	category, err := p.categoryRepo.GetCourseCategory(r.Context(), categoryId)
	if err != nil {
		http.Error(w, "Failed to get course category", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CourseCategoryRes{
		Id:   category.Id,
		Name: category.Name,
	})
}

// ListCourseCategories returns a list of course categories, optionally filtered by search.
// @Summary List course categories
// @Tags Course Category
// @Param search query string false "Search keyword"
// @Success 200 {array} CourseCategoryRes
// @Failure 500 {string} string
// @Security Session
// @Router /courses/categories [get]
func (p *courseCategoryHandler) ListCourseCategories(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.AuthKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	query := r.URL.Query()
	categories, err := p.categoryRepo.ListCourseCategories(r.Context(), query.Get("search"))
	if err != nil {
		http.Error(w, "Failed to get course categories", http.StatusInternalServerError)
		return
	}

	categoryRes := make([]CourseCategoryRes, 0, len(categories))
	for _, category := range categories {
		categoryRes = append(categoryRes, CourseCategoryRes{
			Id:   category.Id,
			Name: category.Name,
		})
	}

	json.NewEncoder(w).Encode(categoryRes)
}
