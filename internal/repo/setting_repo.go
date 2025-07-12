package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"

	"gorm.io/gorm"
)

type SettingRepo interface {
	Create(ctx context.Context, setting *entity.Setting) (uint, error)
	Update(ctx context.Context, setting *entity.Setting) error
	GetByID(ctx context.Context, id uint) (*entity.Setting, error)
}

func NewSettingRepo(
	db *gorm.DB,
) SettingRepo {
	return &SettingRepoImpl{
		db: db,
	}
}

type SettingRepoImpl struct {
	db *gorm.DB
}

// Update implements SettingRepo.
func (s *SettingRepoImpl) Update(ctx context.Context, setting *entity.Setting) error {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := s.db.WithContext(ctx).Save(setting).Error; err != nil {
		return fmt.Errorf("unable to update setting with id %d: %w", setting.ID, common.ErrInternalServerError)
	}

	return nil
}

func (s *SettingRepoImpl) Create(ctx context.Context, setting *entity.Setting) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := s.db.WithContext(ctx).Create(setting).Error; err != nil {
		return 0, fmt.Errorf("unable to create new setting: %w", common.ErrInternalServerError)
	}

	return setting.ID, nil
}

func (s *SettingRepoImpl) GetByID(ctx context.Context, id uint) (*entity.Setting, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	setting := &entity.Setting{}
	if err := s.db.WithContext(ctx).First(setting, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("setting: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get setting with id %d: %w", id, common.ErrInternalServerError)
	}

	return setting, nil
}
