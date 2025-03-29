package model

import (
	"database/sql"
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

func ToCourse(c ICourse) Course {
	return Course{
		Id:          c.Id,
		CategoryId:  c.CategoryId,
		Title:       c.Title,
		Description: c.Description.String,
		Uri:         c.Uri,
		Course:      c.Course,
	}
}

type CourseFilters struct {
	Name       string
	Course     string
	CategoryId string
}
