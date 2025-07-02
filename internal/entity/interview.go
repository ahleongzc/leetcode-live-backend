package entity

type Interview struct {
	Base
	UserID           uint
	QuestionID       uint
	ReviewID         uint
	StartTimestampMS int64
	EndTimestampMS   *int64
	Token            *string
}
