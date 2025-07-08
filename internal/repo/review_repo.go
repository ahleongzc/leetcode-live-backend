package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"

	"gorm.io/gorm"
)

type ReviewRepo interface {
	Create(ctx context.Context, review *entity.Review) (uint, error)
	Update(ctx context.Context, review *entity.Review) error
	GetByID(ctx context.Context, id uint) (*entity.Review, error)
}

func NewReviewRepo(
	db *gorm.DB,
) ReviewRepo {
	return &ReviewRepoImpl{
		db: db,
	}
}

type ReviewRepoImpl struct {
	db *gorm.DB
}

func (r *ReviewRepoImpl) Update(ctx context.Context, review *entity.Review) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := r.db.WithContext(ctx).Save(review).Error; err != nil {
		return fmt.Errorf("unable to update review with id %d: %w", review.ID, common.ErrInternalServerError)
	}

	return nil
}

// Create implements ReviewRepo.
func (r *ReviewRepoImpl) Create(ctx context.Context, review *entity.Review) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		return 0, fmt.Errorf("unable to create new review, %s: %w", err, common.ErrInternalServerError)
	}

	return review.ID, nil
}

// GetByID implements ReviewRepo.
func (r *ReviewRepoImpl) GetByID(ctx context.Context, id uint) (*entity.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	review := &entity.Review{}
	if err := r.db.WithContext(ctx).First(review, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get review with id %d: %w", id, common.ErrInternalServerError)
	}

	return review, nil
}
