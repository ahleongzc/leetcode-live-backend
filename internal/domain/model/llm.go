package model

import "time"

type LLMRole string

const (
	SYSTEM    LLMRole = "system"
	USER      LLMRole = "user"
	ASSISTANT LLMRole = "assistant"
)

type ChatCompletionsRequest struct {
	Messages []*LLMMessage
}

func NewChatCompletionsRequest() *ChatCompletionsRequest {
	return &ChatCompletionsRequest{
		Messages: make([]*LLMMessage, 0),
	}
}

func (c *ChatCompletionsRequest) SetMessages(messages []*LLMMessage) *ChatCompletionsRequest {
	if c == nil {
		return nil
	}
	c.Messages = append(c.Messages, messages...)
	return c
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

func NewLLMMessage() *LLMMessage {
	return &LLMMessage{}
}

func (l *LLMMessage) SetRole(role LLMRole) *LLMMessage {
	if l == nil {
		return nil
	}
	l.Role = role
	return l
}

func (l *LLMMessage) SetContent(content string) *LLMMessage {
	if l == nil {
		return nil
	}
	l.Content = content
	return l
}

type ChatCompletionsResponse struct {
	CreatedTimestampMS int64
	Choices            []*Choice
}

func NewChatCompletionsResponse() *ChatCompletionsResponse {
	return &ChatCompletionsResponse{
		CreatedTimestampMS: time.Now().UnixMilli(),
		Choices:            make([]*Choice, 0),
	}
}

func (c *ChatCompletionsResponse) AppendChoice(choice *Choice) *ChatCompletionsResponse {
	if c == nil || choice == nil {
		return nil
	}

	if c.Choices == nil {
		c.Choices = make([]*Choice, 0)
	}

	c.Choices = append(c.Choices, choice)
	return c
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

func NewChoice() *Choice {
	return &Choice{}
}

func (c *Choice) SetIndex(index int) *Choice {
	if c == nil {
		return nil
	}
	c.Index = index
	return c
}

func (c *Choice) SetMessage(message *LLMMessage) *Choice {
	if c == nil {
		return nil
	}
	c.Message = message
	return c
}

func (c *Choice) GetMessage() *LLMMessage {
	if c == nil {
		return nil
	}
	return c.Message
}
