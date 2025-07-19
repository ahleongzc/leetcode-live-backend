package model

import "github.com/ahleongzc/leetcode-live-backend/internal/util"

type InterviewerResponse struct {
	URL string
	End bool
}

func NewInterviewerResponse() *InterviewerResponse {
	return &InterviewerResponse{}
}

func (i *InterviewerResponse) SetURL(url string) *InterviewerResponse {
	if i == nil {
		return nil
	}
	i.URL = url
	return i
}

func (i *InterviewerResponse) EndInterview() {
	if i == nil {
		return
	}
	i.End = true
}

func (i *InterviewerResponse) Exists() bool {
	return i != nil
}

type InterviewHistory struct {
	Interviews []*Interview `json:"interviews"`
}

func NewInterviewHistory() *InterviewHistory {
	return &InterviewHistory{
		Interviews: make([]*Interview, 0),
	}
}

func (i *InterviewHistory) SetInterviews(interviews []*Interview) *InterviewHistory {
	if i == nil {
		return nil
	}
	i.Interviews = append([]*Interview{}, interviews...)
	return i
}

type Interview struct {
	// This field uses the UUID of the interview for display purposes
	ID string `json:"id"`
	// TODO: This field currently uses the external question ID as the question field, need to see how to change this in the future
	Question             string  `json:"question"`
	QuestionAttemptCount uint    `json:"question_attempt_count"`
	Score                *uint   `json:"score"`
	Passed               *bool   `json:"passed"`
	Feedback             *string `json:"feedback"`
	StartTimestampS      *int64  `json:"start_timestamp_s"`
	EndTimestampS        *int64  `json:"end_timestamp_s"`
	TimeRemainingS       *uint   `json:"time_remaining_s"`
}

func NewInterview() *Interview {
	return &Interview{}
}

func (i *Interview) SetTimeRemainingS(seconds uint) *Interview {
	if i == nil {
		return nil
	}
	i.TimeRemainingS = util.ToPtr(seconds)
	return i
}

// Pass in the UUID here, never use internal id for display
func (i *Interview) SetID(id string) *Interview {
	if i == nil {
		return nil
	}
	i.ID = id
	return i
}

func (i *Interview) SetQuestionAttemptCount(count uint) *Interview {
	if i == nil {
		return nil
	}
	i.QuestionAttemptCount = count
	return i
}

func (i *Interview) SetQuestion(question string) *Interview {
	if i == nil {
		return nil
	}
	i.Question = question
	return i
}

func (i *Interview) SetStartTimestampS(timestampSeconds int64) *Interview {
	if i == nil {
		return nil
	}
	i.StartTimestampS = util.ToPtr(timestampSeconds)
	return i
}

func (i *Interview) SetEndTimestampS(timestampSeconds int64) *Interview {
	if i == nil {
		return nil
	}
	i.EndTimestampS = util.ToPtr(timestampSeconds)
	return i
}

func (i *Interview) SetFeedback(feedback string) *Interview {
	if i == nil {
		return nil
	}
	i.Feedback = util.ToPtr(feedback)
	return i
}

func (i *Interview) SetScore(score uint) *Interview {
	if i == nil {
		return nil
	}
	i.Score = util.ToPtr(score)
	return i
}

func (i *Interview) SetPassed(passed bool) *Interview {
	if i == nil {
		return nil
	}
	i.Passed = util.ToPtr(passed)
	return i
}
