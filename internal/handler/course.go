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

	json.NewEncoder(w).Encode(iCourse.ToCourse())
}

func (c *courseHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value(config.SessionKey)
	if sessionValue == nil {
		c.log.Error("missing session middleware")
	}

	query := r.URL.Query()
	queryFilters := model.CourseFilters{
		Name:       query.Get("search"),
		Course:     query.Get("search"),
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
		courses = append(courses, iCourse.ToCourse())
	}

	json.NewEncoder(w).Encode(courses)
}
