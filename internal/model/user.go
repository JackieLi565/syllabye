package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id          string `json:"id"`
	ProgramId   string `json:"programId,omitempty"`
	FullName    string `json:"fullname"`
	Nickname    string `json:"nickname,omitempty"`
	CurrentYear *int16 `json:"currentYear,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Email       string `json:"email"`
	Picture     string `json:"picture,omitempty"`
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

func (u IUser) ToUser() User {
	return User{
		Id:        u.Id,
		ProgramId: u.ProgramId.String,
		FullName:  u.FullName,
		Nickname:  u.Nickname.String,
		Gender:    u.Gender.String,
		Email:     u.Gender.String,
		Picture:   u.Picture.String,
	}
}

type UserSignUpRequest struct {
	ProgramId   string `json:"program"`
	Nickname    string `json:"nickname"`
	CurrentYear int16  `json:"currentYear"`
	Gender      string `json:"gender"`
}
