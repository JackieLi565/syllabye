package model

import (
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

func (p IProgram) ToProgram() Program {
	return Program{
		Id:        p.Id,
		FacultyId: p.FacultyId,
		Name:      p.Name,
		Uri:       p.Uri,
	}
}

type ProgramFilters struct {
	FacultyId string
	Name      string
}
