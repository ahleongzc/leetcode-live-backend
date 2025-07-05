package entity

import (
	"github.com/ahleongzc/leetcode-live-backend/internal/infra/llm"
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
	URL         string
}

func (t *Transcript) ToLLMMessage() *llm.LLMMessage {
	return &llm.LLMMessage{
		Role:    llm.LLMRole(t.Role),
		Content: t.Content,
	}
}
