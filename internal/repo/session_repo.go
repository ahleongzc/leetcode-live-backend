package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
)

type SessionRepo interface {
	Create(ctx context.Context, session *entity.Session) error
	Update(ctx context.Context, session *entity.Session) error
	GetByID(ctx context.Context, ID string) (*entity.Session, error)
	DeleteByID(ctx context.Context, ID string) error
	DeleteExpired(ctx context.Context) error
}

func NewSessionRepo(
	db *sql.DB,
) SessionRepo {
	return &SessionRepoImpl{
		db: db,
	}
}

type SessionRepoImpl struct {
	db *sql.DB
}

// DeleteExpired implements SessionRepo.
func (s *SessionRepoImpl) DeleteExpired(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{time.Now().UnixMilli()}

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE $1 > expire_timestamp_ms
    `, common.SESSION_TABLE_NAME)

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to delete expired sessions, %s: %w", err.Error(), common.ErrInternalServerError)
	}

	return nil
}

// Update implements SessionRepo.
func (s *SessionRepoImpl) Update(ctx context.Context, session *entity.Session) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{session.UserID, session.ExpireTimestampMS, session.ID}

	query := fmt.Sprintf(`
        UPDATE %s 
        SET user_id = $1, expire_timestamp_ms = $2
        WHERE id = $3
    `, common.SESSION_TABLE_NAME)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to update session with id %s, %s: %w", session.ID, err.Error(), common.ErrInternalServerError)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected when updating session with id %s, %s: %w", session.ID, err.Error(), common.ErrInternalServerError)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("unable to update session with id %s: %w", session.ID, common.ErrInternalServerError)
	}

	return nil
}

// DeleteByID implements SessionRepo.
func (s *SessionRepoImpl) DeleteByID(ctx context.Context, ID string) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{ID}

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, common.SESSION_TABLE_NAME)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to delete session with id %s, %s: %w", ID, err.Error(), common.ErrInternalServerError)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected when deleting session with id %s, %s: %w", ID, err.Error(), common.ErrInternalServerError)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("unable to delete session with id %s: %w", ID, common.ErrNotFound)
	}

	return nil
}

// Create implements SessionRepo.
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

// GetByID implements SessionRepo.
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
