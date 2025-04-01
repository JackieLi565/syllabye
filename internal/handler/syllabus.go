package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/bucket"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/go-chi/chi/v5"
)

type syllabusHandler struct {
	log          logger.Logger
	syllabusRepo repository.SyllabusRepository
	presigner    bucket.PresignerClient
}

func NewSyllabusHandler(log logger.Logger, syllabus repository.SyllabusRepository, presigner bucket.PresignerClient) *syllabusHandler {
	return &syllabusHandler{
		log:          log,
		syllabusRepo: syllabus,
		presigner:    presigner,
	}
}

// GetSyllabus retrieves a specific syllabus by ID and returns a signed URL in the header.
// @Summary Get a syllabus
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Success 200 {object} model.Syllabus
// @Header 200 {string} X-Presigned-Url "Presigned URL to access the syllabus file"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId} [get]
func (s *syllabusHandler) GetSyllabus(w http.ResponseWriter, r *http.Request) {
	sessionValue, ok := r.Context().Value(config.SessionKey).(model.Session)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	syllabusId := chi.URLParam(r, "syllabusId")
	syllabus, err := s.syllabusRepo.GetAndViewSyllabus(r.Context(), sessionValue.UserId, syllabusId)
	if err != nil {
		if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Invalid syllabus ID value.", http.StatusBadRequest)
		} else if errors.Is(err, util.ErrNotFound) {
			http.Error(w, "Syllabus not found.", http.StatusNotFound)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	signedUrl, err := s.presigner.GetObject(r.Context(), syllabus.Id, 60*60)
	if err != nil {
		http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		return
	}

	w.Header().Add("X-Presigned-Url", signedUrl)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.ToSyllabus(syllabus))
}

// CreateSyllabus creates a new syllabus and returns a presigned upload URL in the response header.
// @Summary Create a syllabus
// @Tags Syllabus
// @Accept json
// @Param body body model.CreateSyllabus true "Syllabus data"
// @Success 201 {string} string
// @Header 201 {string} X-Presigned-Url "Presigned URL to upload the syllabus file"
// @Header 201 {string} Location "URL to access the created syllabus"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi [post]
func (s *syllabusHandler) CreateSyllabus(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.SessionKey).(model.Session)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	var body model.CreateSyllabus
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	syllabusId, err := s.syllabusRepo.CreateSyllabus(r.Context(), model.TSyllabus{
		UserId:      session.UserId,
		CourseId:    body.CourseId,
		File:        body.File,
		FileSize:    body.FileSize,
		ContentType: body.ContentType,
		Year:        body.Year,
		Semester:    body.Semester,
	})
	if err != nil {
		if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Invalid body parameter.", http.StatusBadRequest)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	signedUrl, err := s.presigner.PutObject(r.Context(), syllabusId, body.ContentType, body.Checksum, 60*60)
	if err != nil {
		http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		return
	}

	w.Header().Add("X-Presigned-Url", signedUrl)
	w.Header().Set("Location", os.Getenv(config.Domain)+"/syllabi/"+syllabusId)
	w.WriteHeader(http.StatusCreated)
}

// ListSyllabi returns a paginated list of syllabi with optional filters.
// @Summary List syllabi
// @Tags Syllabus
// @Param userId query string false "Filter by user ID"
// @Param courseId query string false "Filter by course ID"
// @Param year query int false "Filter by year"
// @Param semester query string false "Filter by semester"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 10)"
// @Success 200 {array} model.Syllabus
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi [get]
func (s *syllabusHandler) ListSyllabi(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.SessionKey).(model.Session)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()

	var year *int16
	yearQuery, err := strconv.Atoi(query.Get("year"))
	if err == nil {
		yearInt16 := int16(yearQuery)
		year = &yearInt16
	}

	syllabi, err := s.syllabusRepo.ListSyllabi(r.Context(), session.UserId, model.SyllabusFilters{
		UserId:   query.Get("userId"),
		CourseId: query.Get("courseId"),
		Year:     year,
		Semester: query.Get("semester"),
	}, util.NewPaginate(query.Get("page"), query.Get("size")))
	if err != nil {
		if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Invalid user or course ID.", http.StatusBadRequest)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	publicSyllabi := make([]model.Syllabus, 0, len(syllabi))
	for _, syllabus := range syllabi {
		publicSyllabi = append(publicSyllabi, model.ToSyllabus(syllabus))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(publicSyllabi)
}

// UpdateSyllabus updates a syllabus' metadata (year and semester).
// @Summary Update a syllabus
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Param body body model.UpdateSyllabus true "Updated syllabus data"
// @Success 204 {string} string
// @Header 204 {string} Location "URL to access the updated syllabus"
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId} [patch]
func (s *syllabusHandler) UpdateSyllabus(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.SessionKey).(model.Session)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	var body model.UpdateSyllabus
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	syllabusId := chi.URLParam(r, "syllabusId")
	err := s.syllabusRepo.UpdateSyllabus(r.Context(), session.UserId, syllabusId, model.TSyllabus{
		Year:     body.Year,
		Semester: body.Semester,
	})
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			http.Error(w, "Syllabus not found.", http.StatusNotFound)
		} else if errors.Is(err, util.ErrForbidden) {
			http.Error(w, "You do not have access to update this syllabus.", http.StatusForbidden)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Location", os.Getenv(config.Domain)+"/syllabi/"+syllabusId)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteSyllabus removes a syllabus by ID.
// @Summary Delete a syllabus
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Success 204 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId} [delete]
func (s *syllabusHandler) DeleteSyllabus(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.SessionKey).(model.Session)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	syllabusId := chi.URLParam(r, "syllabusId")
	err := s.syllabusRepo.DeleteSyllabus(r.Context(), session.UserId, syllabusId)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			http.Error(w, "Syllabus not found.", http.StatusNotFound)
		} else if errors.Is(err, util.ErrForbidden) {
			http.Error(w, "You do not have access to delete this syllabus.", http.StatusForbidden)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SyllabusReaction registers a reaction to a syllabus.
// @Summary React to a syllabus
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Param body body model.SyllabusReaction true "Reaction action"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 409 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId}/reaction [post]
func (s *syllabusHandler) SyllabusReaction(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.SessionKey).(model.Session)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	var body model.SyllabusReaction
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	syllabusId := chi.URLParam(r, "syllabusId")
	err := s.syllabusRepo.LikeSyllabus(r.Context(), session.UserId, syllabusId, body.Action == "dislike")
	if err != nil {
		if errors.Is(err, util.ErrConflict) {
			http.Error(w, "A conflict has occurred.", http.StatusConflict)
		} else if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Invalid syllabus ID.", http.StatusBadRequest)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteSyllabusReaction removes a user's reaction to a syllabus.
// @Summary Remove syllabus reaction
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId}/reaction [delete]
func (s *syllabusHandler) DeleteSyllabusReaction(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.SessionKey).(model.Session)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	syllabusId := chi.URLParam(r, "syllabusId")
	err := s.syllabusRepo.DeleteSyllabusLike(r.Context(), session.UserId, syllabusId)
	if err != nil {
		if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Invalid syllabus ID.", http.StatusBadRequest)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListSyllabusLikes returns a list of users reactions to a syllabus.
// @Summary List syllabus reactions
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Success 200 {array} model.SyllabusLike
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId}/reactions [get]
func (s *syllabusHandler) ListSyllabusLikes(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	syllabusId := chi.URLParam(r, "syllabusId")
	likes, err := s.syllabusRepo.ListSyllabusLikes(r.Context(), syllabusId)
	if err != nil {
		if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Invalid syllabus ID.", http.StatusBadRequest)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	publicLikes := make([]model.SyllabusLike, 0, len(likes))
	for _, like := range likes {
		publicLikes = append(publicLikes, model.ToSyllabusLike(like))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(publicLikes)
}

// SyncSyllabus triggers synchronization of a syllabus resource.
// @Summary Sync a syllabus
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Success 204 {string} string
// @Header 204 {string} Location "URL to access the synced syllabus"
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Security AWS
// @Router /syllabi/{syllabusId}/sync [get]
func (s *syllabusHandler) SyncSyllabus(w http.ResponseWriter, r *http.Request) {
	// TODO AWS JWT Auth

	syllabusId := chi.URLParam(r, "syllabusId")
	err := s.syllabusRepo.SyncSyllabus(r.Context(), syllabusId)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			http.Error(w, "Syllabus not found.", http.StatusNotFound)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Location", os.Getenv(config.Domain)+"/syllabi/"+syllabusId)
	w.WriteHeader(http.StatusNoContent)
}
