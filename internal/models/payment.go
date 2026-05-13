package models

type Payment struct {
	Id        int
	StudentId int
	Amount    float64
	Status    string
}
