package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Course struct {
	Id          string `json:"id"`
	CategoryId  string `json:"categoryId"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Uri         string `json:"uri"`
	Course      string `json:"course"`
}

type ICourse struct {
	Id          string
	CategoryId  string
	Title       string
	Description sql.NullString
	Uri         string
	Course      string
	DateAdded   time.Time
}

func NewJsonCourse(course ICourse) []byte {
	json, _ := json.Marshal(Course{
		Id:          course.Id,
		CategoryId:  course.CategoryId,
		Title:       course.Title,
		Description: course.Description.String,
		Uri:         course.Uri,
		Course:      course.Course,
	})

	return json
}

type CourseFilters struct {
	Name       string
	Course     string
	CategoryId string
}
