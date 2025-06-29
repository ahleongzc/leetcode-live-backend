package config

import (
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

const (
	LLM_DEV_PROVIDER string = common.OLLAMA
	LLM_DEV_MODEL    string = "gemma3:1b"
	LLM_DEV_BASE_URL string = "http://localhost:11434"
)

type LLMConfig struct {
	Provider string
	Model    string
	BaseURL  string
	APIKey   string
}

func LoadLLMConfig() (*LLMConfig, error) {
	provider := util.GetEnvOr(common.LLM_PROVIDER_KEY, LLM_DEV_PROVIDER)
	model := util.GetEnvOr(common.LLM_MODEL_KEY, LLM_DEV_MODEL)
	baseURL := util.GetEnvOr(common.LLM_BASE_URL_KEY, LLM_DEV_BASE_URL)
	apiKey := util.GetEnvOr(common.LLM_API_KEY, "")

	if provider == "" || model == "" || baseURL == "" {
		return nil, fmt.Errorf("missing llm config, provider=%s model=%s baseURL=%s: %w", provider, model, baseURL, common.ErrInternalServerError)
	}

	if provider != LLM_DEV_PROVIDER && apiKey == "" {
		return nil, fmt.Errorf("missing api key for provider=%s: %w", provider, common.ErrInternalServerError)
	}

	return &LLMConfig{
		Provider: provider,
		Model:    model,
		BaseURL:  baseURL,
		APIKey:   apiKey,
	}, nil
}
