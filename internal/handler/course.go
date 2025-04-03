package handler

import (
	"encoding/json"
	"net/http"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/model"
	"github.com/JackieLi565/syllabye/internal/repository"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/go-chi/chi/v5"
)

type courseHandler struct {
	log        logger.Logger
	courseRepo repository.CourseRepository
}

func NewCourseHandler(log logger.Logger, course repository.CourseRepository) *courseHandler {
	return &courseHandler{
		log:        log,
		courseRepo: course,
	}
}

// GetCourse retrieves a specific course by ID.
// @Summary Get a course
// @Tags Course
// @Param courseId path string true "Course ID"
// @Success 200 {object} model.Course
// @Failure 500 {string} string
// @Security Session
// @Router /courses/{courseId} [get]
func (c *courseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		c.log.Error("missing session middleware")
	}

	courseId := chi.URLParam(r, "courseId")
	iCourse, err := c.courseRepo.GetCourse(r.Context(), courseId)
	if err != nil {
		http.Error(w, "Failed to get course", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.ToCourse(iCourse))
}

// ListCourses returns a paginated list of courses, optionally filtered by name or category.
// @Summary List courses
// @Tags Course
// @Param search query string false "Search by course name or code"
// @Param category query string false "Filter by category ID"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 25)"
// @Success 200 {array} model.Course
// @Failure 500 {string} string
// @Security Session
// @Router /courses [get]
func (c *courseHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		c.log.Error("missing session middleware")
	}

	query := r.URL.Query()
	queryFilters := model.CourseFilters{
		Search:     query.Get("search"),
		CategoryId: query.Get("category"),
	}
	paginateOptions := util.NewPaginate(query.Get("page"), query.Get("size"))
	iCourses, err := c.courseRepo.ListCourses(r.Context(), queryFilters, paginateOptions)
	if err != nil {
		http.Error(w, "Failed to get faculties", http.StatusInternalServerError)
		return
	}

	courses := make([]model.Course, 0, len(iCourses))
	for _, iCourse := range iCourses {
		courses = append(courses, model.ToCourse(iCourse))
	}

	json.NewEncoder(w).Encode(courses)
}
