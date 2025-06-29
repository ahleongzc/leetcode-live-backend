package scenario

import (
	"context"
	"errors"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type QuestionScenario interface {
	// Returns the internal question ID
	GetOrCreateQuestion(ctx context.Context, externalID, description string) (int, error)
}

func NewQuestionScenario(
	questionRepo repo.QuestionRepo,
) QuestionScenario {
	return &QuestionScenarioImpl{
		questionRepo: questionRepo,
	}
}

type QuestionScenarioImpl struct {
	questionRepo repo.QuestionRepo
}

func (q *QuestionScenarioImpl) GetOrCreateQuestion(ctx context.Context, externalID string, description string) (int, error) {
	question, err := q.questionRepo.GetByExternalID(ctx, externalID)
	if nil == err {
		return question.ID, nil
	}

	if !errors.Is(err, common.ErrNotFound) {
		return 0, err
	}

	newQuestion := &entity.Question{
		ExternalID:  externalID,
		Description: description,
	}

	id, err := q.questionRepo.Create(ctx, newQuestion)
	if err != nil {
		return 0, err
	}

	return id, nil
}
