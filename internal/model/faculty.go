package model

import (
	"encoding/json"
	"time"
)

type Faculty struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type IFaculty struct {
	Id        string
	Name      string
	DateAdded time.Time
}

func (f IFaculty) ToFaculty() Faculty {
	return Faculty{
		Id:   f.Id,
		Name: f.Name,
	}
}

func NewJsonFaculty(faculty IFaculty) []byte {
	json, _ := json.Marshal(Faculty{
		Id:   faculty.Id,
		Name: faculty.Name,
	})
	// error can never occur as faculty is an internal model from the database.

	return json
}
