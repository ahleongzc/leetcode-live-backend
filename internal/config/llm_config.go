package config

import (
	"os"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

type LLMModel string

const (
	GPT_4O_MINI LLMModel = "gpt-4o-mini"
)

type LLMConfig struct {
	Model   LLMModel
	BaseURL string
	APIKey  string
}

func LoadLLMConfig() *LLMConfig {
	apiKey := os.Getenv(common.OPENAI_API_KEY)

	return &LLMConfig{
		BaseURL: common.OPENAI_BASE_URL,
		Model:   GPT_4O_MINI,
		APIKey:  apiKey,
	}
}
