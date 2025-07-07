package model

type InterviewHistory struct {
	Interviews []*Interview `json:"interviews"`
}

type Interview struct {
	// This field uses external ID internally
	ID                    string  `json:"id"`
	Question              string  `json:"question"`
	QuestionAttemptNumber uint    `json:"question_attempt_number"`
	Score                 *uint   `json:"score"`
	Passed                *bool   `json:"passed"`
	Feedback              *string `json:"feedback"`
	StartTimestampS       *int64  `json:"start_timestamp_s"`
	EndTimestampS         *int64  `json:"end_timestamp_s"`
}
