package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"

	"gorm.io/gorm"
)

type SessionRepo interface {
	Create(ctx context.Context, session *entity.Session) error
	Update(ctx context.Context, session *entity.Session) error
	GetByToken(ctx context.Context, token string) (*entity.Session, error)
	GetByID(ctx context.Context, id uint) (*entity.Session, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) (uint, error)
}

func NewSessionRepo(
	db *gorm.DB,
) SessionRepo {
	return &SessionRepoImpl{
		db: db,
	}
}

type SessionRepoImpl struct {
	db *gorm.DB
}

// GetByID implements SessionRepo.
func (s *SessionRepoImpl) GetByID(ctx context.Context, id uint) (*entity.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	session := &entity.Session{}
	if err := s.db.WithContext(ctx).First(session, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("session: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get session with id %d: %w", id, common.ErrInternalServerError)
	}

	return session, nil
}

// DeleteExpired implements SessionRepo.
func (s *SessionRepoImpl) DeleteExpired(ctx context.Context) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	result := s.db.WithContext(ctx).
		Where("expire_timestamp_ms < ? ", time.Now().UnixMilli()).
		Delete(&entity.Session{})

	if err := result.Error; err != nil {
		return 0, fmt.Errorf("unable to delete expired session: %w", err)
	}

	return uint(result.RowsAffected), nil
}

// Update implements SessionRepo.
func (s *SessionRepoImpl) Update(ctx context.Context, session *entity.Session) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := s.db.WithContext(ctx).Save(session).Error; err != nil {
		return fmt.Errorf("unable to update session: %w", common.ErrInternalServerError)
	}

	return nil
}

// DeleteByToken implements SessionRepo.
func (s *SessionRepoImpl) DeleteByToken(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	result := s.db.WithContext(ctx).
		Where("token = ?", token).
		Delete(&entity.Session{})

	if err := result.Error; err != nil {
		return fmt.Errorf("unable to delete session with token %s, %s: %w", token, err, common.ErrInternalServerError)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found %s: %w", token, common.ErrNotFound)
	}

	return nil
}

// Create implements SessionRepo.
func (s *SessionRepoImpl) Create(ctx context.Context, session *entity.Session) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("unable to create new session, %s: %w", err, common.ErrInternalServerError)
	}

	return nil
}

// GetByID implements SessionRepo.
func (s *SessionRepoImpl) GetByToken(ctx context.Context, token string) (*entity.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	var session entity.Session
	err := s.db.WithContext(ctx).First(&session, "token = ?", token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("session: %w", common.ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("unable to get session with token %s, %s: %w", token, err, common.ErrInternalServerError)
	}
	return &session, nil
}
