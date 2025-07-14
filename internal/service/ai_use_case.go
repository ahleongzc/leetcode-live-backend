package service

import (
	"context"
	"io"

	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type AIUseCase interface {
	GenerateSpeechReply(ctx context.Context, text, instruction string) (io.Reader, error)
	GenerateTextReply(ctx context.Context, messages []*model.LLMMessage) (string, error)
}

func NewAIUseCase(
	ttsRepo repo.TTSRepo,
	llmRepo repo.LLMRepo,
) AIUseCase {
	return &AIUseCaseImpl{
		ttsRepo: ttsRepo,
		llmRepo: llmRepo,
	}
}

type AIUseCaseImpl struct {
	ttsRepo repo.TTSRepo
	llmRepo repo.LLMRepo
}

// GenerateSpeechReply implements AIService.
func (a *AIUseCaseImpl) GenerateSpeechReply(ctx context.Context, text, instruction string) (io.Reader, error) {
	reader, err := a.ttsRepo.TextToSpeechReader(ctx, text, instruction)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

// GenerateTextReply implements AIService.
func (a *AIUseCaseImpl) GenerateTextReply(ctx context.Context, messages []*model.LLMMessage) (string, error) {
	req := model.NewChatCompletionsRequest().
		SetMessages(messages)

	resp, err := a.llmRepo.ChatCompletions(ctx, req)
	if err != nil {
		return "", err
	}

	reply := resp.GetResponse().GetContent()
	return reply, nil
}
