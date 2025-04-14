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
	"github.com/oapi-codegen/nullable"
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

type UserResponse struct {
	Id          string                    `json:"id"`
	ProgramId   nullable.Nullable[string] `json:"programId,omitempty" swaggertype:"primitive,string" extensions:"x-nullable"`
	FullName    string                    `json:"fullname,omitempty"`
	Nickname    nullable.Nullable[string] `json:"nickname,omitempty" swaggertype:"primitive,string" extensions:"x-nullable"`
	CurrentYear nullable.Nullable[int16]  `json:"currentYear,omitempty" swaggertype:"primitive,string" extensions:"x-nullable"`
	Gender      nullable.Nullable[string] `json:"gender,omitempty" swaggertype:"primitive,string" extensions:"x-nullable"`
	Email       string                    `json:"email,omitempty"`
	Picture     nullable.Nullable[string] `json:"picture,omitempty" swaggertype:"primitive,string" extensions:"x-nullable"`
	Bio         nullable.Nullable[string] `json:"bio,omitempty" swaggertype:"primitive,string" extensions:"x-nullable"`
} //@name UserResponse

// GetUser retrieves a user by ID.
// @Summary Get a user
// @Tags User
// @Param userId path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /users/{userId} [get]
func (u *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.AuthKey)
	if sessionValue == nil {
		u.log.Error("missing session middleware")
	}

	userId := chi.URLParam(r, "userId")
	user, err := u.userRepo.GetUser(r.Context(), userId)
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
	json.NewEncoder(w).Encode(UserResponse{
		Id:          user.Id,
		ProgramId:   util.DefaultNullable(user.ProgramId.Valid, user.ProgramId.String),
		FullName:    user.FullName,
		Nickname:    util.DefaultNullable(user.Nickname.Valid, user.Nickname.String),
		CurrentYear: util.DefaultNullable(user.CurrentYear.Valid, user.CurrentYear.Int16),
		Gender:      util.DefaultNullable(user.Gender.Valid, user.Gender.String),
		Email:       user.Email,
		Picture:     util.DefaultNullable(user.Picture.Valid, user.Picture.String),
		Bio:         util.DefaultNullable(user.Bio.Valid, user.Bio.String),
	})
}

type UpdateUserRequest struct {
	ProgramId   nullable.Nullable[string] `json:"programId" swaggertype:"primitive,string" extensions:"x-nullable"`
	Nickname    nullable.Nullable[string] `json:"nickname" swaggertype:"primitive,string" extensions:"x-nullable"`
	CurrentYear nullable.Nullable[int16]  `json:"currentYear" swaggertype:"primitive,integer" extensions:"x-nullable"`
	Gender      nullable.Nullable[string] `json:"gender" swaggertype:"primitive,string" extensions:"x-nullable"`
	Bio         nullable.Nullable[string] `json:"bio" swaggertype:"primitive,string" extensions:"x-nullable"`
} //@name UpdateUserRequest

// UpdateUser modifies a user's profile data.
// @Summary Update a user
// @Tags User
// @Param userId path string true "User ID"
// @Param body body UpdateUserRequest true "Updated user data"
// @Success 201 {string} string
// @Header 201 {string} Location "URL to access the updated user"
// @Failure 400 {string} string
// @Failure 403 {string} string
// @Failure 409 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /users/{userId} [patch]
func (u *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.AuthKey).(model.Session)
	if !ok {
		u.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	userId := chi.URLParam(r, "userId")
	// A user can only update themselves (based on the same user id of session)
	if userId != session.UserId {
		http.Error(w, "You're not allowed to modify another user's profile.", http.StatusForbidden)
		return
	}

	var body UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := u.userRepo.UpdateUser(r.Context(), userId, repository.UpdateUser{
		ProgramId:   body.ProgramId,
		Nickname:    body.Nickname,
		CurrentYear: body.CurrentYear,
		Gender:      body.Gender,
		Bio:         body.Bio,
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

	w.Header().Set("Location", os.Getenv(config.ServerDomain)+"/users/"+userId)
	w.WriteHeader(http.StatusCreated)
}

// AddUserCourse adds a course to a user's academic history.
// @Summary Add a user course
// @Tags User
// @Param userId path string true "User ID"
// @Param body body model.CreateUserCourse true "User course data"
// @Success 201 {string} string
// @Failure 400 {string} string
// @Failure 403 {string} string
// @Failure 409 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /users/{userId}/courses [post]
func (u *userHandler) AddUserCourse(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.AuthKey).(model.Session)
	if !ok {
		u.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	userId := chi.URLParam(r, "userId")
	if userId != session.UserId {
		http.Error(w, "You're not allowed to modify another user's profile.", http.StatusForbidden)
		return
	}

	var body model.CreateUserCourse
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	if body.CourseId == "" {
		http.Error(w, "Invalid or missing request body fields.", http.StatusBadRequest)
		return
	}

	err := u.userRepo.AddUserCourse(r.Context(), session.UserId, model.TUserCourse{
		CourseId:      body.CourseId,
		YearTaken:     body.YearTaken,
		SemesterTaken: body.SemesterTaken,
	})
	if err != nil {
		if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Malformed request data.", http.StatusBadRequest)
		} else if errors.Is(err, util.ErrConflict) {
			http.Error(w, "A configuration error occurred.", http.StatusConflict)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteUserCourse removes a course from a user's academic history.
// @Summary Delete a user course
// @Tags User
// @Param userId path string true "User ID"
// @Param courseId path string true "Course ID"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /users/{userId}/courses/{courseId} [delete]
func (u *userHandler) DeleteUserCourse(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.AuthKey).(model.Session)
	if !ok {
		u.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	userId := chi.URLParam(r, "userId")
	if userId != session.UserId {
		http.Error(w, "You're not allowed to modify another user's profile.", http.StatusForbidden)
		return
	}

	err := u.userRepo.DeleteUserCourse(r.Context(), session.UserId, chi.URLParam(r, "courseId"))
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			http.Error(w, "Course not found.", http.StatusNotFound)
		} else if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Malformed request data.", http.StatusBadRequest)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateUserCourse updates a user's course information.
// @Summary Update a user course
// @Tags User
// @Param userId path string true "User ID"
// @Param courseId path string true "Course ID"
// @Param body body model.UpdateUserCourse true "Updated course data"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Failure 409 {string} string
// @Failure 500 {string} string
// @Security Session
// @Router /users/{userId}/courses/{courseId} [patch]
func (u *userHandler) UpdateUserCourse(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(config.AuthKey).(model.Session)
	if !ok {
		u.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	userId := chi.URLParam(r, "userId")
	if userId != session.UserId {
		http.Error(w, "You're not allowed to modify another user's profile.", http.StatusForbidden)
		return
	}

	var body model.UpdateUserCourse
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	if body.SemesterTaken == "" || body.YearTaken == 0 {
		http.Error(w, "Invalid or missing request body fields.", http.StatusBadRequest)
		return
	}

	err := u.userRepo.UpdateUserCourse(r.Context(), session.UserId, chi.URLParam(r, "courseId"), model.TUserCourse{
		YearTaken:     &body.YearTaken,
		SemesterTaken: &body.SemesterTaken,
	})
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			http.Error(w, "Course not found.", http.StatusNotFound)
		} else if errors.Is(err, util.ErrMalformed) {
			http.Error(w, "Malformed request data.", http.StatusBadRequest)
		} else if errors.Is(err, util.ErrConflict) {
			http.Error(w, "A configuration error occurred.", http.StatusConflict)
		} else {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListUserCourses retrieves a paginated list of a user's courses.
// @Summary List user courses
// @Tags User
// @Param userId path string true "User ID"
// @Param search query string false "Search by name or course code"
// @Param category query string false "Filter by category ID"
// @Param page query string false "Page number (default: 1)"
// @Param size query string false "Page size (default: 25)"
// @Success 200 {array} model.UserCourse
// @Failure 500 {string} string
// @Security Session
// @Router /users/{userId}/courses [get]
func (u *userHandler) ListUserCourses(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.AuthKey)
	if sessionValue == nil {
		u.log.Error("session middleware potential missing")
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()
	queryFilters := model.CourseFilters{
		Search:     query.Get("search"),
		CategoryId: query.Get("category"),
	}

	courses, err := u.userRepo.ListUserCourses(
		r.Context(),
		chi.URLParam(r, "userId"),
		queryFilters,
		util.NewPaginate(query.Get("page"), query.Get("size")),
	)
	if err != nil {
		http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		return
	}

	userCourses := make([]model.UserCourse, 0, len(courses))
	for _, course := range courses {
		userCourses = append(userCourses, model.UserCourse{
			CourseId:      course.CourseId,
			Title:         course.Title,
			Course:        course.Course,
			YearTaken:     course.YearTaken.Int16,
			SemesterTaken: course.SemesterTaken.String,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userCourses)
}

// SearchUserNickname checks if a user nickname exists.
// @Summary Check existing nickname
// @Tags User
// @Param search query string false "Search user nickname"
// @Success 200 {object} model.UserNicknameExists
// @Failure 500 {string} string
// @Security Session
// @Router /users/exists [get]
func (u *userHandler) SearchUserNickname(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(model.UserNicknameExists{
			Exists: false,
		})
		return
	}

	exists, err := u.userRepo.SearchUserNickname(r.Context(), search)
	if err != nil {
		http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.UserNicknameExists{
		Exists: exists,
	})
}
