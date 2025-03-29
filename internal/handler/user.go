package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/go-chi/chi/v5"
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

func (u *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		u.log.Error("missing session middleware")
	}

	userId := chi.URLParam(r, "userId")
	iUser, err := u.userRepo.GetUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			http.Error(w, "User not found.", http.StatusNotFound)
		} else if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
		} else {
			u.log.Error("unhandled error when getting user", logger.Err(err))
			http.Error(w, "An unknown error occurred.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.ToUser(iUser))
}

func (u *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		u.log.Error("missing session middleware")
	}
	session, ok := sessionValue.(model.ISession)
	if !ok {
		u.log.Warn("session type asset failed")
		http.Error(w, "Should always be protected", http.StatusUnauthorized)
		return
	}

	userId := chi.URLParam(r, "userId")
	// A user can only update themselves (based on the same user id of session)
	if userId != session.UserId {
		http.Error(w, "You're not allowed to update other peoples profile!", http.StatusForbidden)
		return
	}

	var body model.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := u.userRepo.UpdateUser(r.Context(), userId, model.TUser{
		ProgramId:   body.ProgramId,
		Nickname:    body.Nickname,
		CurrentYear: body.CurrentYear,
		Gender:      body.Gender,
	})
	if err != nil {
		if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		} else if errors.Is(err, util.ErrConflict) {
			http.Error(w, "A conflict has occurred", http.StatusConflict)
		} else {
			http.Error(w, "A internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Location", os.Getenv(config.Domain)+"/users/"+userId)
	w.WriteHeader(http.StatusCreated)
}
