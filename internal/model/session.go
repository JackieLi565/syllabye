package model

import (
	"encoding/json"
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

func NewJsonSession(session ISession) []byte {
	json, _ := json.Marshal(Session{
		UserId:  session.UserId,
		Expired: session.DateExpires.Before(time.Now()),
	})

	return json
}
