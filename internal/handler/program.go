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

type programHandler struct {
	log         logger.Logger
	programRepo repository.ProgramRepository
}

func NewProgramHandler(log logger.Logger, program repository.ProgramRepository) *programHandler {
	return &programHandler{
		log:         log,
		programRepo: program,
	}
}

func (p *programHandler) GetProgram(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	programId := chi.URLParam(r, "programId")
	iProgram, err := p.programRepo.GetProgram(r.Context(), programId)
	if err != nil {
		http.Error(w, "Failed to get program", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(iProgram.ToProgram())
}

func (p *programHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	query := r.URL.Query()
	iPrograms, err := p.programRepo.ListPrograms(r.Context(), model.ProgramFilters{
		FacultyId: query.Get("faculty"),
		Name:      query.Get("search"),
	})
	if err != nil {
		http.Error(w, "Failed to get programs", http.StatusInternalServerError)
		return
	}

	programs := make([]model.Program, 0, len(iPrograms))
	for _, iProgram := range iPrograms {
		programs = append(programs, iProgram.ToProgram())
	}

	json.NewEncoder(w).Encode(programs)
}
