package entity

type Interview struct {
	ID                 int
	UserID             int
	ExternalQuestionID string
	StartTimestampMS   int64
	EndTimestampMS     *int64
}
