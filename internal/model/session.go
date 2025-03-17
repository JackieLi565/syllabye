package model

import "time"

type Session struct {
	Id          string
	UserId      string
	DateAdded   time.Time
	DateExpires time.Time
	User        *User
}

type SessionResponse struct {
	UserId string `json:"userId"`
}
