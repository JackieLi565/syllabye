package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type User struct {
	Id           string
	ProgramId    sql.NullString
	FullName     string
	Nickname     sql.NullString
	CurrentYear  sql.NullInt16
	Gender       sql.NullString
	Email        string
	Picture      sql.NullString
	IsActive     bool
	DateAdded    time.Time
	DateModified time.Time
}

type UserResponse struct {
	Id          string          `json:"id"`
	Program     json.RawMessage `json:"program"`
	FullName    string          `json:"fullName"`
	Nickname    *string         `json:"nickname"`
	CurrentYear *int16          `json:"currentYear"`
	Gender      *string         `json:"gender"`
	Email       string          `json:"email"`
	Picture     *string         `json:"picture"`
}
