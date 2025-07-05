package llm

import (
	"context"
	"net/http"
)

type OpenAI struct {
	model      string
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewOpenAILLM(model, baseURL, apiKey string, httpClient *http.Client) LLM {
	return &OpenAI{
		model:      model,
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (o *OpenAI) ChatCompletions(ctx context.Context, chatCompletionsRequest *ChatCompletionsRequest) (*ChatCompletionsResponse, error) {
	panic("implement me")
}
