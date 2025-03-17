package model

import (
	"encoding/json"
	"time"
)

type Program struct {
	Id        string
	FacultyId string
	Faculty   *Faculty
	Name      string
	URI       string
	DateAdded time.Time
}

type ProgramResponse struct {
	Id      string          `json:"id"`
	Faculty json.RawMessage `json:"faculty"`
	Name    string          `json:"name"`
	Uri     string          `json:"uri"`
}
