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
