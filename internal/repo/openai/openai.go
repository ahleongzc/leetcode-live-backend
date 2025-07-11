package openai

import (
	"context"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
)

type OpenAILLM struct {
	model      string
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewOpenAILLM(model, baseURL, apiKey string, httpClient *http.Client) *OpenAILLM {
	return &OpenAILLM{
		model:      model,
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (o *OpenAILLM) ChatCompletions(ctx context.Context, chatCompletionsRequest *model.ChatCompletionsRequest) (*model.ChatCompletionsResponse, error) {
	panic("implement me")
}
