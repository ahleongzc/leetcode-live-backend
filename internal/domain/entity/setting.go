package entity

import "time"

type Setting struct {
	Base
	InterviewDurationS      uint
	RemainingInterviewCount uint
}

func NewDefaultSetting() *Setting {
	setting := NewSetting()
	setting.SetInterviewDurationS(uint(20 * time.Minute / time.Second))
	setting.SetRemainingInterviewCount(20)

	return setting
}

func NewSetting() *Setting {
	return &Setting{}
}

func (s *Setting) SetInterviewDurationS(durationSeconds uint) *Setting {
	if s == nil {
		return nil
	}
	s.InterviewDurationS = durationSeconds
	return s
}

func (s *Setting) SetRemainingInterviewCount(count uint) *Setting {
	if s == nil {
		return s
	}
	s.RemainingInterviewCount = count
	return s
}

func (s *Setting) GetRemainingInterviewCount() uint {
	if s == nil {
		return 0
	}
	return s.RemainingInterviewCount
}

func (s *Setting) DecrementRemainingInterviewCount() {
	if s == nil {
		return
	}

	remainingCount := s.GetRemainingInterviewCount()
	if remainingCount == 0 {
		return
	}

	s.SetRemainingInterviewCount(remainingCount - 1)
}
