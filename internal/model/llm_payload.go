package model

import (
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

type LLMRole string

const (
	SYSTEM    LLMRole = "system"
	USER      LLMRole = "user"
	ASSISTANT LLMRole = "assistant"
)

type ChatCompletionsRequest struct {
	Messages []*LLMMessage
}

func (c *ChatCompletionsRequest) GetMessages() []*LLMMessage {
	if c == nil {
		return nil
	}
	if c.Messages == nil {
		c.Messages = make([]*LLMMessage, 0)
	}
	return c.Messages
}

type LLMMessage struct {
	Role    LLMRole
	Content string
}

func (l *LLMMessage) GetRole() LLMRole {
	if l == nil {
		return SYSTEM
	}
	return l.Role
}

func (l *LLMMessage) GetContent() string {
	if l == nil {
		return ""
	}
	return l.Content
}

func NewMessage(role LLMRole, content string) *LLMMessage {
	return &LLMMessage{
		Role:    role,
		Content: content,
	}
}

type ChatCompletionsResponse struct {
	CreatedTimestampMS int64
	Choices            []*Choice
}

func NewChatCompletionsResponse() *ChatCompletionsResponse {
	return &ChatCompletionsResponse{
		Choices: make([]*Choice, 0),
	}
}

func (c *ChatCompletionsResponse) AddChoice(choice *Choice) error {
	if c == nil {
		return fmt.Errorf("unable to add choice to nil chatCompletionsResponseModel %w", common.ErrInternalServerError)
	}
	if choice == nil {
		return fmt.Errorf("trying to add nil choice to chatCompletionsResponseModel %w", common.ErrInternalServerError)
	}

	if c.Choices == nil {
		c.Choices = make([]*Choice, 0)
	}

	c.Choices = append(c.Choices, choice)
	return nil
}

func (c *ChatCompletionsResponse) GetResponse() *LLMMessage {
	if c == nil {
		return nil
	}
	if len(c.Choices) == 0 {
		return nil
	}

	return c.Choices[0].GetMessage()
}

type Choice struct {
	Index   int
	Message *LLMMessage
}

func NewChoice(index int, message *LLMMessage) *Choice {
	return &Choice{
		Index:   index,
		Message: message,
	}
}

func (c *Choice) GetMessage() *LLMMessage {
	if c == nil {
		return nil
	}

	return c.Message
}
