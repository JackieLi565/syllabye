package model

import (
	"time"
)

type Session struct {
	UserId  string `json:"userId"`
	Expired bool   `json:"expired"`
}

type ISession struct {
	Id          string
	UserId      string
	DateAdded   time.Time
	DateExpires time.Time
}

func (s ISession) ToSession() Session {
	return Session{
		UserId:  s.UserId,
		Expired: s.DateExpires.Before(time.Now()),
	}
}
