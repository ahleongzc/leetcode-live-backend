package entity

type Interview struct {
	Base
	UserID           uint
	QuestionID       uint
	StartTimestampMS int64
	ReviewID         *uint
	EndTimestampMS   *int64
	Token            *string
}
