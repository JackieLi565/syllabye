package model

import (
	"database/sql"
	"time"
)

type UserCourse struct {
	CourseId      string `json:"courseId"`
	Title         string `json:"title"`
	Course        string `json:"course"`
	YearTaken     int16  `json:"yearTaken,omitempty"`
	SemesterTaken string `json:"semesterTaken,omitempty"`
}

type IUserCourse struct {
	UserId        string
	CourseId      string
	Title         string
	Course        string
	YearTaken     sql.NullInt16
	SemesterTaken sql.NullString
	DateAdded     time.Time
	DateModified  time.Time
}

type TUserCourse struct {
	UserId        string
	CourseId      string
	YearTaken     *int16
	SemesterTaken *string
	DateAdded     time.Time
	DateModified  time.Time
}

type UpdateUserCourse struct {
	YearTaken     int16  `json:"yearTaken"`
	SemesterTaken string `json:"semesterTaken"`
}

type CreateUserCourse struct {
	YearTaken     *int16  `json:"yearTaken"`
	SemesterTaken *string `json:"semesterTaken"`
	CourseId      string  `json:"courseId"`
}

type UserNicknameExists struct {
	Exists bool `json:"exists"`
}
