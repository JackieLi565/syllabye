package handler

import (
	"encoding/json"
	"net/http"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/go-chi/chi/v5"
)

type facultyHandler struct {
	log         logger.Logger
	facultyRepo repository.FacultyRepository
}

func NewFacultyHandler(log logger.Logger, faculty repository.FacultyRepository) *facultyHandler {
	return &facultyHandler{
		log:         log,
		facultyRepo: faculty,
	}
}

type FacultyRes struct {
	Id   string `json:"id"`
	Name string `json:"name"`
} //@name FacultyResponse

// GetFaculty retrieves a specific faculty by ID.
// @Summary Get a faculty
// @Tags Faculty
// @Param facultyId path string true "Faculty ID"
// @Success 200 {object} FacultyResponse
// @Failure 500 {string} string
// @Security Session
// @Router /faculties/{facultyId} [get]
func (p *facultyHandler) GetFaculty(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.AuthKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	facultyId := chi.URLParam(r, "facultyId")
	faculty, err := p.facultyRepo.GetFaculty(r.Context(), facultyId)
	if err != nil {
		http.Error(w, "Failed to get faculty", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(FacultyRes{
		Id:   faculty.Id,
		Name: faculty.Name,
	})
}

// ListFaculties returns a list of faculties, optionally filtered by search keyword.
// @Summary List faculties
// @Tags Faculty
// @Param search query string false "Search by faculty name"
// @Success 200 {array} FacultyResponse
// @Failure 500 {string} string
// @Security Session
// @Router /faculties [get]
func (p *facultyHandler) ListFaculties(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.AuthKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	query := r.URL.Query()
	faculties, err := p.facultyRepo.ListFaculties(r.Context(), query.Get("search"))
	if err != nil {
		http.Error(w, "Failed to get faculties", http.StatusInternalServerError)
		return
	}

	facultyRes := make([]FacultyRes, 0, len(faculties))
	for _, faculty := range faculties {
		facultyRes = append(facultyRes, FacultyRes{
			Id:   faculty.Id,
			Name: faculty.Name,
		})
	}

	json.NewEncoder(w).Encode(facultyRes)
}
