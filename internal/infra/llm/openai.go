package llm

import (
	"context"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/model"
)

type OpenAI struct {
	model      string
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewOpenAILLM(model, baseURL, apiKey string, httpClient *http.Client) *OpenAI {
	return &OpenAI{
		model:      model,
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (o *OpenAI) ChatCompletions(ctx context.Context, chatCompletionsRequest *model.ChatCompletionsRequest) (*model.ChatCompletionsResponse, error) {
	panic("implement me")
}
