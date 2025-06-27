package model

import (
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

type ChatCompletionsRequest struct {
	Messages []*Message
}

func (c *ChatCompletionsRequest) GetMessages() []*Message {
	if c == nil {
		return nil
	}
	if c.Messages == nil {
		c.Messages = make([]*Message, 0)
	}
	return c.Messages
}

type Message struct {
	Role    string
	Content string
}

func NewMessage(role, content string) *Message {
	return &Message{
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

func (c *ChatCompletionsResponse) GetResponse() *Message {
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
	Message *Message
}

func NewChoice(index int, message *Message) *Choice {
	return &Choice{
		Index:   index,
		Message: message,
	}
}

func (c *Choice) GetMessage() *Message {
	if c == nil {
		return nil
	}

	return c.Message
}
