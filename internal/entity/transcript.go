package entity

import "github.com/ahleongzc/leetcode-live-backend/internal/model"

type Role string

const (
	SYSTEM    Role = "system"
	USER      Role = "user"
	ASSISTANT Role = "assistant"
)

type Transcript struct {
	ID                 int
	Role               Role
	Content            string
	InterviewID        int
	CreatedTimestampMS int64
	URL                string
}

func (t *Transcript) ToLLMMessage() *model.LLMMessage {
	return &model.LLMMessage{
		Role:    model.LLMRole(t.Role),
		Content: t.Content,
	}
}
