package model

import (
	"time"
)

type Faculty struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type IFaculty struct {
	Id        string
	Name      string
	DateAdded time.Time
}

func ToFaculty(f IFaculty) Faculty {
	return Faculty{
		Id:   f.Id,
		Name: f.Name,
	}
}
