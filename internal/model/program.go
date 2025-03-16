package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Program struct {
	ID        pgtype.UUID `json:"id"`
	FacultyID pgtype.UUID `json:"faculty_id"`
	Name      string      `json:"name"`
	URI       string      `json:"uri"`
	DateAdded time.Time   `json:"date_added"`
}
