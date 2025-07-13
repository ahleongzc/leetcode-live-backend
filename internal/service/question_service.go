package service

import (
	"context"
	"errors"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type QuestionService interface {
	// Returns the internal question ID
	GetOrCreateQuestion(ctx context.Context, externalID, description string) (uint, error)
}

func NewQuestionService(
	questionRepo repo.QuestionRepo,
) QuestionService {
	return &QuestionServiceImpl{
		questionRepo: questionRepo,
	}
}

type QuestionServiceImpl struct {
	questionRepo repo.QuestionRepo
}

func (q *QuestionServiceImpl) GetOrCreateQuestion(ctx context.Context, externalID string, description string) (uint, error) {
	question, err := q.questionRepo.GetByExternalID(ctx, externalID)
	if nil == err {
		return question.ID, nil
	}

	if !errors.Is(err, common.ErrNotFound) {
		return 0, err
	}

	newQuestion := entity.NewQuestion().
		SetExternalID(externalID).
		SetDescription(description)

	id, err := q.questionRepo.Create(ctx, newQuestion)
	if err != nil {
		return 0, err
	}

	return id, nil
}
