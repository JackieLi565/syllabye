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

func (p *facultyHandler) GetFaculty(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	facultyId := chi.URLParam(r, "facultyId")
	iFaculty, err := p.facultyRepo.GetFaculty(r.Context(), facultyId)
	if err != nil {
		http.Error(w, "Failed to get faculty", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(iFaculty.ToFaculty())
}

func (p *facultyHandler) ListFaculties(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	query := r.URL.Query()
	iFaculties, err := p.facultyRepo.ListFaculties(r.Context(), query.Get("search"))
	if err != nil {
		http.Error(w, "Failed to get faculties", http.StatusInternalServerError)
		return
	}

	faculties := make([]model.Faculty, 0, len(iFaculties))
	for _, iFaculty := range iFaculties {
		faculties = append(faculties, iFaculty.ToFaculty())
	}

	json.NewEncoder(w).Encode(faculties)
}
