package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id          string `json:"id"`
	ProgramId   string `json:"programId,omitempty"`
	FullName    string `json:"fullname,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	CurrentYear *int16 `json:"currentYear,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Email       string `json:"email,omitempty"`
	Picture     string `json:"picture,omitempty"`
}

func ToUser(u IUser) User {
	return User{
		Id:        u.Id,
		ProgramId: u.ProgramId.String,
		FullName:  u.FullName,
		Nickname:  u.Nickname.String,
		Gender:    u.Gender.String,
		Email:     u.Email,
		Picture:   u.Picture.String,
	}
}

type IUser struct {
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

type TUser struct {
	Id           string
	ProgramId    string
	FullName     string
	Nickname     string
	CurrentYear  int16
	Gender       string
	Email        string
	Picture      string
	IsActive     bool
	DateAdded    time.Time
	DateModified time.Time
}

type UpdateUser struct {
	ProgramId   string `json:"program"`
	Nickname    string `json:"nickname"`
	CurrentYear int16  `json:"currentYear"`
	Gender      string `json:"gender"`
}
