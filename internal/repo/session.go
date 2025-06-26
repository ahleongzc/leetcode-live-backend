package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
)

type SessionRepo interface {
	Create(ctx context.Context, session *entity.Session) error
	GetByID(ctx context.Context, ID string) (*entity.Session, error)
	DeleteByID(ctx context.Context, ID string) error
}

func NewSessionRepoImpl(
	db *sql.DB,
) SessionRepo {
	return &SessionRepoImpl{
		db: db,
	}
}

type SessionRepoImpl struct {
	db *sql.DB
}

// DeleteSessionByID implements SessionRepo.
func (s *SessionRepoImpl) DeleteByID(ctx context.Context, ID string) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{ID}

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, common.SESSION_TABLE_NAME)

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to delete session with id %s, %s: %w", ID, err.Error(), common.ErrInternalServerError)
	}

	return nil
}

// CreateSession implements SessionRepo.
func (s *SessionRepoImpl) Create(ctx context.Context, session *entity.Session) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{session.ID, session.ExpireTimestampMS}

	query := fmt.Sprintf(`
		INSERT INTO %s
		(id, expire_timestamp_ms)
		VALUES
		($1, $2)
	`, common.SESSION_TABLE_NAME)

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to create new session, %s: %w", err.Error(), common.ErrInternalServerError)
	}

	return nil
}

// GetSessionByID implements SessionRepo.
func (s *SessionRepoImpl) GetByID(ctx context.Context, ID string) (*entity.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{ID}

	query := fmt.Sprintf(`
		SELECT 
			id, expire_timestamp_ms 
		FROM 
			%s
		WHERE 
			id = $1
	`, common.SESSION_TABLE_NAME)

	session := &entity.Session{}
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&session.ID, &session.ExpireTimestampMS)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get session with id %s, %s: %w", ID, err.Error(), common.ErrInternalServerError)
	}

	return session, nil
}
