package entity

import "time"

type Setting struct {
	Base
	InterviewDurationS      int64
	RemainingInterviewCount uint
}

func NewDefaultSetting() *Setting {
	return &Setting{
		InterviewDurationS:      int64(20 * time.Minute / time.Second),
		RemainingInterviewCount: 20,
	}
}
