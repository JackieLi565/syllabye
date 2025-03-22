package model

import (
	"database/sql"
	"encoding/json"
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

func NewJsonUser(user IUser) []byte {
	jsonUser := User{
		Id:        user.Id,
		ProgramId: user.ProgramId.String,
		FullName:  user.FullName,
		Nickname:  user.Nickname.String,
		Gender:    user.Gender.String,
		Email:     user.Gender.String,
		Picture:   user.Picture.String,
	}

	if user.CurrentYear.Valid {
		jsonUser.CurrentYear = &user.CurrentYear.Int16
	} else {
		jsonUser.CurrentYear = nil
	}

	json, _ := json.Marshal(user)

	return json
}

type UserSignUpRequest struct {
	ProgramId   string `json:"program"`
	Nickname    string `json:"nickname"`
	CurrentYear int16  `json:"currentYear"`
	Gender      string `json:"gender"`
}
