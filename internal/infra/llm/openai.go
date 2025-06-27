package llm

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/model"
)

type OpenAI struct {
	Model string
}

func NewOpenAILLM(model string) *OpenAI {
	return &OpenAI{
		Model: model,
	}
}

func (o *OpenAI) ChatCompletions(ctx context.Context, chatCompletionsRequest *model.ChatCompletionsRequest) (*model.ChatCompletionsResponse, error) {
	panic("implement me")
}
