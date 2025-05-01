package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/authorizer"
	"github.com/JackieLi565/syllabye/internal/service/bucket"
	"github.com/JackieLi565/syllabye/internal/service/emailer"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/service/queue"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/nullable"
)

type syllabusHandler struct {
	log          logger.Logger
	syllabusRepo repository.SyllabusRepository
	presigner    bucket.PresignerClient
	jwt          *authorizer.JwtAuthorizer
	queue        queue.WebhookQueue
	emailer      emailer.NoReplyEmailer
}

func NewSyllabusHandler(log logger.Logger, syllabus repository.SyllabusRepository, presigner bucket.PresignerClient, jwt *authorizer.JwtAuthorizer, queue queue.WebhookQueue, emailer emailer.NoReplyEmailer) *syllabusHandler {
	return &syllabusHandler{
		log:          log,
		syllabusRepo: syllabus,
		presigner:    presigner,
		jwt:          jwt,
		queue:        queue,
		emailer:      emailer,
	}
}

type SyllabusRes struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	CourseId    string `json:"courseId"`
	File        string `json:"fileName"`
	FileSize    int    `json:"fileSize"`
	ContentType string `json:"contentType"`
	Year        int16  `json:"year"`
	Semester    string `json:"semester"`
	DateAdded   int64  `json:"dateAdded"`
	Received    bool   `json:"received"`
} //@name SyllabusResponse

// GetSyllabus retrieves a specific syllabus by ID and returns a signed URL in the header.
// @Summary Get a syllabus
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Success 200 {object} SyllabusResponse
// @Header 200 {string} X-Presigned-Url "Presigned URL to access the syllabus file"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId} [get]
func (s *syllabusHandler) GetSyllabus(w http.ResponseWriter, r *http.Request) {
	sessionValue, ok := r.Context().Value(config.AuthKey).(SessionPayload)
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
	json.NewEncoder(w).Encode(SyllabusRes{
		Id:          syllabus.Id,
		UserId:      syllabus.UserId,
		CourseId:    syllabus.CourseId,
		File:        syllabus.File,
		FileSize:    syllabus.FileSize,
		ContentType: syllabus.ContentType,
		Year:        syllabus.Year,
		Semester:    syllabus.Semester,
		DateAdded:   syllabus.DateAdded.UnixMicro(),
		Received:    syllabus.DateSynced.Valid,
	})
}

type AddSyllabusReq struct {
	CourseId    string `json:"courseId" validate:"required"`
	File        string `json:"fileName" validate:"required"`
	FileSize    int    `json:"fileSize" validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
	Checksum    string `json:"checksum" validate:"required"`
	Year        int16  `json:"year" validate:"required"`
	Semester    string `json:"semester" validate:"required"`
} //@name CreateSyllabusRequest

// CreateSyllabus creates a new syllabus and returns a presigned upload URL in the response header.
// @Summary Create a syllabus
// @Tags Syllabus
// @Accept json
// @Param body body CreateSyllabusRequest true "Syllabus data"
// @Success 201 {string} string
// @Header 201 {string} X-Presigned-Url "Presigned URL to upload the syllabus file"
// @Header 201 {string} Location "URL to access the created syllabus"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi [post]
func (s *syllabusHandler) CreateSyllabus(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.AuthKey).(SessionPayload)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	var body AddSyllabusReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	syllabusId, err := s.syllabusRepo.CreateSyllabus(r.Context(), repository.InsertSyllabus{
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

	var delaySeconds int32
	if os.Getenv(config.ENV) == "development" {
		delaySeconds = 10
	} else {
		delaySeconds = 60 * 5 // 5 minutes
	}

	// Clean up syllabus API
	token, err := s.jwt.EncodeJwt(nil)
	if err != nil {
		http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
	}
	requestId, _ := r.Context().Value(config.RequestIdKey).(string)
	s.queue.SendMessage(r.Context(), queue.WebhookMessage{
		RequestId: requestId,
		Url:       os.Getenv(config.ServerDomain) + "/syllabi/" + syllabusId + "/verify",
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	}, delaySeconds)

	w.Header().Add("X-Presigned-Url", signedUrl)
	w.Header().Set("Location", os.Getenv(config.ServerDomain)+"/syllabi/"+syllabusId)
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
// @Success 200 {array} SyllabusResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi [get]
func (s *syllabusHandler) ListSyllabi(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.AuthKey).(SessionPayload)
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

	syllabi, err := s.syllabusRepo.ListSyllabi(r.Context(), session.UserId, repository.SyllabusFilters{
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

	publicSyllabi := make([]SyllabusRes, 0, len(syllabi))
	for _, syllabus := range syllabi {
		publicSyllabi = append(publicSyllabi, SyllabusRes{
			Id:          syllabus.Id,
			UserId:      syllabus.UserId,
			CourseId:    syllabus.CourseId,
			File:        syllabus.File,
			FileSize:    syllabus.FileSize,
			ContentType: syllabus.ContentType,
			Year:        syllabus.Year,
			Semester:    syllabus.Semester,
			DateAdded:   syllabus.DateAdded.UnixMicro(),
			Received:    syllabus.DateSynced.Valid,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(publicSyllabi)
}

type UpdateSyllabusReq struct {
	Year     nullable.Nullable[int16]  `json:"year" swaggertype:"primitive,integer" extensions:"x-nullable"`
	Semester nullable.Nullable[string] `json:"semester" swaggertype:"primitive,string" extensions:"x-nullable"`
} //@name UpdateSyllabusRequest

// UpdateSyllabus updates a syllabus' metadata (year and semester).
// @Summary Update a syllabus
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Param body body UpdateSyllabusRequest true "Updated syllabus data"
// @Success 204 {string} string
// @Header 204 {string} Location "URL to access the updated syllabus"
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId} [patch]
func (s *syllabusHandler) UpdateSyllabus(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.AuthKey).(SessionPayload)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	var body UpdateSyllabusReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	syllabusId := chi.URLParam(r, "syllabusId")
	err := s.syllabusRepo.UpdateSyllabus(r.Context(), session.UserId, syllabusId, repository.UpdateSyllabus{
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

	w.Header().Add("Location", os.Getenv(config.ServerDomain)+"/syllabi/"+syllabusId)
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
	session, ok := r.Context().Value(config.AuthKey).(SessionPayload)
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

type SyllabusLikeReq struct {
	Action string `json:"action"`
} //@name SyllabusReactionRequest

// SyllabusReaction registers a reaction to a syllabus.
// @Summary React to a syllabus
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Param body body SyllabusReactionRequest true "Reaction action"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 409 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId}/reaction [post]
func (s *syllabusHandler) SyllabusReaction(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.AuthKey).(SessionPayload)
	if !ok {
		s.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	var body SyllabusLikeReq
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
	session, ok := r.Context().Value(config.AuthKey).(SessionPayload)
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

type SyllabusLikeRes struct {
	SyllabusId string `json:"syllabusId"`
	UserId     string `json:"userId"`
	IsDislike  bool   `json:"dislike"`
	DateAdded  int64  `json:"dateReacted"`
} //@name SyllabusReactionResponse

// ListSyllabusLikes returns a list of users reactions to a syllabus.
// @Summary List syllabus reactions
// @Tags Syllabus
// @Param syllabusId path string true "Syllabus ID"
// @Success 200 {array} SyllabusReactionResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /syllabi/{syllabusId}/reactions [get]
func (s *syllabusHandler) ListSyllabusLikes(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.AuthKey)
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

	publicLikes := make([]SyllabusLikeRes, 0, len(likes))
	for _, like := range likes {
		publicLikes = append(publicLikes, SyllabusLikeRes{
			SyllabusId: like.SyllabusId,
			UserId:     like.UserId,
			IsDislike:  like.IsDislike,
			DateAdded:  like.DateAdded.UnixMicro(),
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(publicLikes)
}

// SyncSyllabus triggers synchronization of a syllabus resource.
func (s *syllabusHandler) SyncSyllabus(w http.ResponseWriter, r *http.Request) {
	syllabusId := chi.URLParam(r, "syllabusId")
	err := s.syllabusRepo.SyncSyllabus(r.Context(), syllabusId)
	if err != nil {
		http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		return
	}

	_, meta, err := s.syllabusRepo.VerifySyllabus(r.Context(), syllabusId)
	if err != nil {
		http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		return
	}

	s.emailer.SendSubmissionMissingEmail(r.Context(), meta.UserEmail, meta.UserName, meta.Course)

	s.log.Info(fmt.Sprintf("syllabus %s synced", syllabusId))
	w.WriteHeader(http.StatusNoContent)
}

// VerifySyllabus check if a syllabus has been synchronized.
func (s *syllabusHandler) VerifySyllabus(w http.ResponseWriter, r *http.Request) {
	syllabusId := chi.URLParam(r, "syllabusId")
	isVerified, meta, err := s.syllabusRepo.VerifySyllabus(r.Context(), syllabusId)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			http.Error(w, "Syllabus not found.", http.StatusNotFound)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	if !isVerified {
		s.log.Info(fmt.Sprintf("syllabus %s not verified", syllabusId))
		s.emailer.SendSubmissionMissingEmail(r.Context(), meta.UserEmail, meta.UserName, meta.Course)
	} else {
		s.log.Info(fmt.Sprintf("syllabus %s verified", syllabusId))
	}

	w.WriteHeader(http.StatusNoContent)
}
