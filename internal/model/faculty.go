package model

import "time"

type Faculty struct {
	Id        string
	Name      string
	DateAdded time.Time
}

type FacultyResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
