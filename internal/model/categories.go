package model

import "encoding/json"

type CourseCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ICourseCategory struct {
	Id        string
	Name      string
	DateAdded string
}

func NewJsonCourseCategory(category ICourseCategory) []byte {
	json, _ := json.Marshal(CourseCategory{
		Id:   category.Id,
		Name: category.Name,
	})

	return json
}
