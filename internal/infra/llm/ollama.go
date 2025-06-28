package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
)

type Ollama struct {
	model      string
	host       string
	httpClient *http.Client
}

func NewOllamaLLM(model string) *Ollama {
	return &Ollama{
		model: model,
		host:  common.OLLAMA_BASE_URL,
		httpClient: &http.Client{
			Timeout: common.HTTP_REQUEST_TIMEOUT,
		},
	}
}

func (o *Ollama) ChatCompletions(ctx context.Context, chatCompletionsRequest *model.ChatCompletionsRequest) (*model.ChatCompletionsResponse, error) {
	if chatCompletionsRequest == nil {
		return nil, fmt.Errorf("chatCompletionRequest cannot be nil when calling Ollama: %w", common.ErrInternalServerError)
	}

	url := o.host + "/v1/chat/completions"
	ollamaReq, err := o.convertToOllamaChatCompletionsRequest(chatCompletionsRequest)
	if err != nil {
		return nil, err
	}

	jsonPayload, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal payload before calling ollama, %s: %w", err.Error(), common.ErrInternalServerError)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to perform HTTP call to ollama, %s: %w", err.Error(), common.ErrInternalServerError)
	}
	req.Header.Set(common.CONTENT_TYPE, "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform HTTP call to ollama, %s: %w", err.Error(), common.ErrInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response from Ollama is not ok: %w", common.ErrInternalServerError)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to HTTP call to ollama, %s: %w", err.Error(), common.ErrInternalServerError)
	}

	var ollamaResp OllamaChatCompletionsResponse
	err = json.Unmarshal(body, &ollamaResp)
	if err != nil {
		return nil, fmt.Errorf("unable to  HTTP call to ollama, %s: %w", err.Error(), common.ErrInternalServerError)
	}

	chatCompletionsResponseModel, err := o.convertToChatCompletionsResponseModel(&ollamaResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), common.ErrInternalServerError)
	}

	return chatCompletionsResponseModel, nil
}

type OllamaChatCompletionsRequest struct {
	Model    string
	Messages []*OllamaMessage
}

func NewOllamaChatCompletionsRequest() *OllamaChatCompletionsRequest {
	return &OllamaChatCompletionsRequest{
		Messages: make([]*OllamaMessage, 0),
	}
}

type OllamaChatCompletionsResponse struct {
	ID                string          `json:"id"`
	Model             string          `json:"model"`
	CreatedTimestampS int64           `json:"created"`
	Choices           []*OllamaChoice `json:"choices"`
}

func NewOllamaChatCompletionsResponse() *OllamaChatCompletionsResponse {
	return &OllamaChatCompletionsResponse{
		Choices: make([]*OllamaChoice, 0),
	}
}

func (o *OllamaChatCompletionsResponse) getChoices() []*OllamaChoice {
	if o == nil {
		return nil
	}
	if o.Choices == nil {
		return make([]*OllamaChoice, 0)
	}

	return o.Choices
}

type OllamaChoice struct {
	Index   int            `json:"index"`
	Message *OllamaMessage `json:"message"`
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (o *OllamaChatCompletionsRequest) addMessage(message *OllamaMessage) {
	if o.Messages == nil {
		o.Messages = make([]*OllamaMessage, 0)
	}
	o.Messages = append(o.Messages, message)
}

func (o *Ollama) convertToOllamaChatCompletionsRequest(req *model.ChatCompletionsRequest) (*OllamaChatCompletionsRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("cannot convert a nil chatCompletionsRequestModel into ollamaChatCompletionsRequest: %w", common.ErrInternalServerError)
	}

	ollamaChatCompletionsRequest := NewOllamaChatCompletionsRequest()
	ollamaChatCompletionsRequest.Model = o.model

	for _, message := range req.GetMessages() {
		ollamaMessage := &OllamaMessage{
			Content: message.Content,
			Role:    message.Role,
		}
		ollamaChatCompletionsRequest.addMessage(ollamaMessage)
	}

	return ollamaChatCompletionsRequest, nil
}

func (o *OllamaMessage) convertToMessage() (*model.Message, error) {
	if o == nil {
		return nil, fmt.Errorf("cannot convert nil ollamaMessage into message: %w", common.ErrInternalServerError)
	}
	return &model.Message{
		Role:    o.Role,
		Content: o.Content,
	}, nil
}

func (o *OllamaChoice) convertToChoice() (*model.Choice, error) {
	if o == nil {
		return nil, fmt.Errorf("cannot convert nil ollamaChoice into choice: %w", common.ErrInternalServerError)
	}
	message, err := o.Message.convertToMessage()
	if err != nil {
		return nil, err
	}

	return &model.Choice{
		Index:   o.Index,
		Message: message,
	}, nil
}

func (o *Ollama) convertToChatCompletionsResponseModel(ollamaResp *OllamaChatCompletionsResponse) (*model.ChatCompletionsResponse, error) {
	chatCompletionsResponseModel := model.NewChatCompletionsResponse()

	for _, ollamaChoice := range ollamaResp.getChoices() {
		choice, err := ollamaChoice.convertToChoice()
		if err != nil {
			return nil, fmt.Errorf("%s, %w", err.Error(), common.ErrInternalServerError)
		}
		chatCompletionsResponseModel.AddChoice(choice)
	}

	return chatCompletionsResponseModel, nil
}
