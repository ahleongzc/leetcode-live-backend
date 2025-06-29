package infra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra/llm"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
)

type LLM interface {
	ChatCompletions(ctx context.Context, request *model.ChatCompletionsRequest) (*model.ChatCompletionsResponse, error)
}

func NewLLM(
	llmConfig *config.LLMConfig, httpClient *http.Client,
) (LLM, error) {
	switch llmConfig.Provider {
	case config.LLM_DEV_PROVIDER:
		return llm.NewOllamaLLM(llmConfig.Model, llmConfig.BaseURL, httpClient), nil
	case common.OPENAI:
		return llm.NewOpenAILLM(llmConfig.Model, llmConfig.BaseURL, llmConfig.APIKey, httpClient), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider %s: %w", llmConfig.Provider, common.ErrInternalServerError)
	}
}
