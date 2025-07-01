package repo

import (
	"context"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"

	"gorm.io/gorm"
)

type InterviewRepo interface {
	Create(ctx context.Context, interview *entity.Interview) error
	Update(ctx context.Context, interview *entity.Interview) error
	GetByToken(ctx context.Context, token string) (*entity.Interview, error)
	GetByID(ctx context.Context, id uint) (*entity.Interview, error)
}

func NewInterviewRepo(
	db *gorm.DB,
) InterviewRepo {
	return &InterviewRepoImpl{
		db: db,
	}
}

type InterviewRepoImpl struct {
	db *gorm.DB
}

func (i *InterviewRepoImpl) GetByID(ctx context.Context, id uint) (*entity.Interview, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	interview := &entity.Interview{}
	if err := i.db.WithContext(ctx).
		First(interview, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("interview: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get interview with id %d, %s: %w", id, err.Error(), common.ErrInternalServerError)
	}

	return interview, nil
}

func (i *InterviewRepoImpl) GetByToken(ctx context.Context, token string) (*entity.Interview, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	interview := &entity.Interview{}
	if err := i.db.WithContext(ctx).
		Where("token = ?", token).
		First(interview).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("interview: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get interview with token %s, %s: %w", token, err.Error(), common.ErrInternalServerError)
	}

	return interview, nil
}

// Create implements InterviewRepo.
func (i *InterviewRepoImpl) Create(ctx context.Context, interview *entity.Interview) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := i.db.WithContext(ctx).Create(interview).Error; err != nil {
		return fmt.Errorf("unable to create new interview: %w", common.ErrInternalServerError)
	}

	return nil
}

// Update implements InterviewRepo.
func (i *InterviewRepoImpl) Update(ctx context.Context, interview *entity.Interview) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := i.db.WithContext(ctx).Save(interview).Error; err != nil {
		return fmt.Errorf("unable to update interview with id %d: %w", interview.ID, common.ErrInternalServerError)
	}

	return nil
}
