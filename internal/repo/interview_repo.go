package repo

import (
	"context"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"

	"gorm.io/gorm"
)

type InterviewRepo interface {
	Create(ctx context.Context, interview *entity.Interview) (uint, error)
	Update(ctx context.Context, interview *entity.Interview) error
	GetByToken(ctx context.Context, token string) (*entity.Interview, error)
	GetByID(ctx context.Context, id uint) (*entity.Interview, error)
	GetUnfinishedInterviewByUserID(ctx context.Context, userID uint) (*entity.Interview, error)
	GetUnstartedInterviewByUserID(ctx context.Context, userID uint) (*entity.Interview, error)
	GetOngoingInterviewByUserID(ctx context.Context, userID uint) (*entity.Interview, error)
	CountByUserIDAndQuestionID(ctx context.Context, userID, questionID uint) (uint, error)
	ListStartedInterviewsByUserID(ctx context.Context, userID, limit, offset uint) ([]*entity.Interview, uint, error)
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

// GetOngoingInterviewByUserID implements InterviewRepo.
func (i *InterviewRepoImpl) GetOngoingInterviewByUserID(ctx context.Context, userID uint) (*entity.Interview, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	interview := &entity.Interview{}
	if err := i.db.WithContext(ctx).
		Where("user_id = ? AND ongoing IS true", userID).
		First(interview).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("interview: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get unstarted interview for user id %d, %s: %w", userID, err, common.ErrInternalServerError)
	}

	return interview, nil
}

func (i *InterviewRepoImpl) CountByUserIDAndQuestionID(ctx context.Context, userID, questionID uint) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	var count int64

	if err := i.db.WithContext(ctx).
		Model(&entity.Interview{}).
		Where("user_id = ? AND question_id = ?", userID, questionID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("unable to count interviews for user ID %d and question ID %d: %w",
			userID, questionID, common.ErrInternalServerError)
	}

	return uint(count), nil
}

// ListByUserID implements InterviewRepo.
func (i *InterviewRepoImpl) ListStartedInterviewsByUserID(ctx context.Context, userID, limit, offset uint) ([]*entity.Interview, uint, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	var interviews []*entity.Interview
	var total int64

	if err := i.db.WithContext(ctx).
		Model(&entity.Interview{}).
		Where("user_id = ? AND start_timestamp_ms IS NOT NULL", userID).
		Count(&total).
		Error; err != nil {
		return nil, 0, fmt.Errorf("unable to count interviews for user ID %d: %w",
			userID, common.ErrInternalServerError)
	}

	result := i.db.WithContext(ctx).
		Where("user_id = ? AND start_timestamp_ms IS NOT NULL", userID).
		Order("end_timestamp_ms IS NULL DESC").
		Order("end_timestamp_ms DESC").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&interviews)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("unable to list interviews for user ID %d: %w",
			userID, common.ErrInternalServerError)
	}

	return interviews, uint(total), nil
}

// GetOngoingInterviewByUserID implements InterviewRepo.
func (i *InterviewRepoImpl) GetUnstartedInterviewByUserID(ctx context.Context, userID uint) (*entity.Interview, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	interview := &entity.Interview{}
	if err := i.db.WithContext(ctx).
		Where("user_id = ? AND start_timestamp_ms IS NULL", userID).
		First(interview).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("interview: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get unstarted interview for user id %d, %s: %w", userID, err, common.ErrInternalServerError)
	}

	return interview, nil
}

// GetOngoingInterviewByUserID implements InterviewRepo.
func (i *InterviewRepoImpl) GetUnfinishedInterviewByUserID(ctx context.Context, userID uint) (*entity.Interview, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	interview := &entity.Interview{}
	if err := i.db.WithContext(ctx).
		Where("user_id = ? AND start_timestamp_ms IS NOT NULL AND end_timestamp_ms IS NULL", userID).
		First(interview).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("interview: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get unfinished interview for user id %d, %s: %w", userID, err, common.ErrInternalServerError)
	}

	return interview, nil
}

func (i *InterviewRepoImpl) GetByID(ctx context.Context, id uint) (*entity.Interview, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
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
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
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
func (i *InterviewRepoImpl) Create(ctx context.Context, interview *entity.Interview) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := i.db.WithContext(ctx).Create(interview).Error; err != nil {
		return 0, fmt.Errorf("unable to create new interview: %w", common.ErrInternalServerError)
	}

	return interview.ID, nil
}

// Update implements InterviewRepo.
func (i *InterviewRepoImpl) Update(ctx context.Context, interview *entity.Interview) error {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := i.db.WithContext(ctx).Save(interview).Error; err != nil {
		return fmt.Errorf("unable to update interview with id %d: %w", interview.ID, common.ErrInternalServerError)
	}

	return nil
}
