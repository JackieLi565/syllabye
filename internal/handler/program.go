package handler

import (
	"encoding/json"
	"net/http"

	"github.com/JackieLi565/syllabye/internal/config"
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

type ProgramRes struct {
	Id        string `json:"id"`
	FacultyId string `json:"faculty"`
	Name      string `json:"name"`
	Uri       string `json:"uri"`
} //@name ProgramResponse

// GetProgram retrieves a specific program by ID.
// @Summary Get a program
// @Tags Program
// @Param programId path string true "Program ID"
// @Success 200 {object} ProgramResponse
// @Failure 500 {string} string
// @Security Session
// @Router /programs/{programId} [get]
func (p *programHandler) GetProgram(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.AuthKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	programId := chi.URLParam(r, "programId")
	program, err := p.programRepo.GetProgram(r.Context(), programId)
	if err != nil {
		http.Error(w, "Failed to get program", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ProgramRes{
		Id:        program.Id,
		FacultyId: program.FacultyId,
		Name:      program.Name,
		Uri:       program.Uri,
	})
}

// ListPrograms returns a list of programs, optionally filtered by faculty or name.
// @Summary List programs
// @Tags Program
// @Param faculty query string false "Filter by faculty ID"
// @Param search query string false "Search by program name or code"
// @Success 200 {array} ProgramResponse
// @Failure 500 {string} string
// @Security Session
// @Router /programs [get]
func (p *programHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.AuthKey)
	if sessionValue == nil {
		p.log.Error("missing session middleware")
	}

	query := r.URL.Query()
	programs, err := p.programRepo.ListPrograms(r.Context(), repository.ProgramFilters{
		FacultyId: query.Get("faculty"),
		Name:      query.Get("search"),
	})
	if err != nil {
		http.Error(w, "Failed to get programs", http.StatusInternalServerError)
		return
	}

	programRes := make([]ProgramRes, 0, len(programs))
	for _, program := range programs {
		programRes = append(programRes, ProgramRes{
			Id:        program.Id,
			FacultyId: program.FacultyId,
			Name:      program.Name,
			Uri:       program.Uri,
		})
	}

	json.NewEncoder(w).Encode(programRes)
}
