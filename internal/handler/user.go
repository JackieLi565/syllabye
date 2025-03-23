package handler

import (
	"encoding/json"
	"net/http"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/logger"
)

type userHandler struct {
	log      logger.Logger
	userRepo repository.UserRepository
}

func NewUserHandler(log logger.Logger, user repository.UserRepository) *userHandler {
	return &userHandler{
		log:      log,
		userRepo: user,
	}
}

func (u *userHandler) CompleteSignUp(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	session, ok := sessionValue.(*model.Session)
	if !ok || session == nil {
		u.log.Error("missing session middleware")
		http.Error(w, "Should always be protected", http.StatusUnauthorized)
		return
	}

	var body model.UserSignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := u.userRepo.CompleteUserSignUp(session.UserId, body)
	if err != nil {
		http.Error(w, "Failed to complete user sign up", http.StatusInternalServerError)
		return
	}

	// No content or direct to loc of user?
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Success"))
}
