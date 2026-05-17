package models

type Group struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Level     string `json:"level"`
	TeacherId int    `json:"teacher_id"`
}
