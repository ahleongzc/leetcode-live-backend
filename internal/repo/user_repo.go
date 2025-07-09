package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	DeleteByID(ctx context.Context, id uint) error
}

func NewUserRepo(
	db *gorm.DB,
) UserRepo {
	return &UserRepoImpl{
		db: db,
	}
}

type UserRepoImpl struct {
	db *gorm.DB
}

// GetByEmail implements UserRepo.
func (u *UserRepoImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	user := &entity.User{}
	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get user with email %s, %s: %w", email, err, common.ErrInternalServerError)
	}

	return user, nil
}

// Create implements UserRepo.
func (u *UserRepoImpl) Create(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("unable to create new user: %w", common.ErrInternalServerError)
	}

	return nil
}

// Delete implements UserRepo.
func (u *UserRepoImpl) DeleteByID(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	result := u.db.WithContext(ctx).Delete(&entity.User{}, id)
	if err := result.Error; err != nil {
		return fmt.Errorf("unable to delete user with id %d: %w", id, common.ErrInternalServerError)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found: %w", common.ErrNotFound)
	}

	return nil
}

// GetByID implements UserRepo.
func (u *UserRepoImpl) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	user := &entity.User{}
	if err := u.db.WithContext(ctx).First(user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get user with id %d: %w", id, common.ErrInternalServerError)
	}

	return user, nil
}

// Update implements UserRepo.
func (u *UserRepoImpl) Update(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := u.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("unable to update user with id %d: %w", user.ID, common.ErrInternalServerError)
	}

	return nil
}
