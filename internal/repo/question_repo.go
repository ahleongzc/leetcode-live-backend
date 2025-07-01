package repo

import (
	"context"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"

	"gorm.io/gorm"
)

type QuestionRepo interface {
	Create(ctx context.Context, question *entity.Question) (uint, error)
	GetByExternalID(ctx context.Context, externalID string) (*entity.Question, error)
}

func NewQuestionRepo(
	db *gorm.DB,
) QuestionRepo {
	return &QuestionRepoImpl{
		db: db,
	}
}

type QuestionRepoImpl struct {
	db *gorm.DB
}

// Create implements QuestionRepo.
func (q *QuestionRepoImpl) Create(ctx context.Context, question *entity.Question) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := q.db.WithContext(ctx).Create(question).Error; err != nil {
		return 0, fmt.Errorf("unable to create new question: %w", err)
	}

	return question.ID, nil
}

// GetByExternalID implements QuestionRepo.
func (q *QuestionRepoImpl) GetByExternalID(ctx context.Context, externalID string) (*entity.Question, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	question := &entity.Question{}

	if err := q.db.WithContext(ctx).
		Where("external_id = ?", externalID).
		First(question).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("question: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get question with external_id %s: %w", externalID, err)
	}

	return question, nil
}
