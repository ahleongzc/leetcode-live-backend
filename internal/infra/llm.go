package infra

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/infra/llm"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type LLM interface {
	ChatCompletions(ctx context.Context, request *model.ChatCompletionsRequest) (*model.ChatCompletionsResponse, error)
}

func NewLLM() LLM {
	if util.IsDevEnv() {
		return llm.NewOllamaLLM("gemma3:1b")
	}
	if util.IsProdEnv() {
		return llm.NewOpenAILLM("gpt-4o-mini")
	}
	return nil
}
