package entity

type Review struct {
	Base
	Score    uint
	Passed   bool
	Feedback string
}
