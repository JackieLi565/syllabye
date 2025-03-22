package model

import (
	"encoding/json"
	"time"
)

type Program struct {
	Id        string `json:"id"`
	FacultyId string `json:"faculty"`
	Name      string `json:"name"`
	Uri       string `json:"uri"`
}

type IProgram struct {
	Id        string
	FacultyId string
	Name      string
	Uri       string
	DateAdded time.Time
}

func NewJsonProgram(program IProgram) []byte {
	json, _ := json.Marshal(Program{
		Id:        program.Id,
		FacultyId: program.FacultyId,
		Name:      program.Name,
		Uri:       program.Uri,
	})

	return json
}

type ProgramFilters struct {
	FacultyId string
	Name      string
}
