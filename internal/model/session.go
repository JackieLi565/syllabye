package model

import (
	"time"
)

type Session struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

type ISession struct {
	Id        string
	UserId    string
	DateAdded time.Time
}
