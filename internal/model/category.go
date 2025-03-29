package model

import "time"

type CourseCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ICourseCategory struct {
	Id        string
	Name      string
	DateAdded time.Time
}

func ToCourseCategory(c ICourseCategory) CourseCategory {
	return CourseCategory{
		Id:   c.Id,
		Name: c.Name,
	}
}
