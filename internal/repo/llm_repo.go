package repo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/ollama"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/openai"
)

type LLMRepo interface {
	ChatCompletions(ctx context.Context, request *model.ChatCompletionsRequest) (*model.ChatCompletionsResponse, error)
}

func NewLLMRepo(
	llmConfig *config.LLMConfig, httpClient *http.Client,
) (LLMRepo, error) {
	switch llmConfig.Provider {
	case config.LLM_DEV_PROVIDER:
		return ollama.NewOllamaLLM(llmConfig.Model, llmConfig.BaseURL, httpClient), nil
	case common.OPENAI:
		return openai.NewOpenAILLM(llmConfig.Model, llmConfig.BaseURL, llmConfig.APIKey, httpClient), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider %s: %w", llmConfig.Provider, common.ErrInternalServerError)
	}
}
