package entity

type Interview struct {
	ID               int
	UserID           int
	QuestionID       int
	StartTimestampMS int64
	EndTimestampMS   *int64
	Token            *string
}
