package models

import "time"

type Lesson struct {
	Id      int
	GroupId int
	Date    time.Time
	Status  string
}
