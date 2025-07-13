package entity

import (
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type Role string

const (
	SYSTEM    Role = "system"
	USER      Role = "user"
	ASSISTANT Role = "assistant"
)

type Intent string

const (
	NO_INTENT             Intent = ""
	HINT_REQUEST          Intent = "hint"
	CLARIFICATION_REQUEST Intent = "clarification"
	END_REQUEST           Intent = "end"
)

type Transcript struct {
	Base
	Role        Role
	Content     string
	InterviewID uint
	Intent      Intent
	URL         *string
}

func (t *Transcript) ToLLMMessage() *model.LLMMessage {
	return &model.LLMMessage{
		Role:    model.LLMRole(t.Role),
		Content: t.Content,
	}
}

func NewTranscript() *Transcript {
	return &Transcript{}
}

func NewInterviewerTranscript() *Transcript {
	transcript := NewTranscript()
	transcript.SetRole(ASSISTANT)

	return transcript
}

func NewCandidateTranscript() *Transcript {
	transcript := NewTranscript()
	transcript.SetRole(USER)

	return transcript
}

func (t *Transcript) SetRole(role Role) *Transcript {
	if t == nil {
		return nil
	}
	t.Role = role
	return t
}

func (t *Transcript) SetContent(content string) *Transcript {
	if t == nil {
		return nil
	}
	t.Content = content
	return t
}

func (t *Transcript) SetInterviewID(interviewID uint) *Transcript {
	if t == nil {
		return nil
	}
	t.InterviewID = interviewID
	return t
}

func (t *Transcript) SetURL(url string) *Transcript {
	if t == nil {
		return nil
	}
	t.URL = util.ToPtr(url)
	return t
}
