package model

import (
	"database/sql"
	"time"
)

type Syllabus struct {
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
}

type ISyllabus struct {
	Id          string
	UserId      string
	CourseId    string
	File        string
	FileSize    int
	ContentType string
	Year        int16
	Semester    string
	DateAdded   time.Time
	DateSynced  sql.NullTime
}

type TSyllabus struct {
	Id          string
	UserId      string
	CourseId    string
	File        string
	FileSize    int
	ContentType string
	Year        int16
	Semester    string
	DateAdded   time.Time
}

type CreateSyllabus struct {
	CourseId    string `json:"courseId"`
	File        string `json:"fileName"`
	FileSize    int    `json:"fileSize"`
	ContentType string `json:"contentType"`
	Checksum    string `json:"checksum"`
	Year        int16  `json:"year"`
	Semester    string `json:"semester"`
}

type UpdateSyllabus struct {
	Year     int16  `json:"year,omitempty"`
	Semester string `json:"semester,omitempty"`
}

type SyllabusReaction struct {
	Action string `json:"action"`
}

type SyllabusFilters struct {
	UserId   string
	CourseId string
	Year     *int16
	Semester string
}

func ToSyllabus(syllabus ISyllabus) Syllabus {
	return Syllabus{
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
	}
}
