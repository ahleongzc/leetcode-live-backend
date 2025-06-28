package service

import (
	"context"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
)

type InterviewService interface {
	ProcessMessage(ctx context.Context, message *model.InterviewMessage) (*model.InterviewMessage, error)
}

func NewInterviewService(
	llm infra.LLM,
) InterviewService {
	return &InterviewServiceImpl{
		llm: llm,
	}
}

type InterviewServiceImpl struct {
	llm infra.LLM
}

func (i *InterviewServiceImpl) ProcessMessage(ctx context.Context, message *model.InterviewMessage) (*model.InterviewMessage, error) {
	req := &model.ChatCompletionsRequest{
		Messages: []*model.Message{
			{
				Role:    "user",
				Content: "write a short greeting",
			},
		},
	}

	resp, err := i.llm.ChatCompletions(ctx, req)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.GetResponse())

	return nil, nil
}
