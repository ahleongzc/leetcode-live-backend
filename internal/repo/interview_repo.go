package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
)

type InterviewRepo interface {
	Create(ctx context.Context, interview *entity.Interview) error
	Update(ctx context.Context, interview *entity.Interview) error
	GetByToken(ctx context.Context, token string) (*entity.Interview, error)
}

func NewInterviewRepo(
	db *sql.DB,
) InterviewRepo {
	return &InterviewRepoImpl{
		db: db,
	}
}

type InterviewRepoImpl struct {
	db *sql.DB
}

func (i *InterviewRepoImpl) GetByToken(ctx context.Context, token string) (*entity.Interview, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{token}

	query := fmt.Sprintf(`
		SELECT 
			id, user_id, question_id, start_timestamp_ms, end_timestamp_ms, token
		FROM 
			%s
		WHERE 
			token = $1
	`, common.INTERVIEW_TABLE_NAME)

	interview := &entity.Interview{}
	err := i.db.QueryRowContext(ctx, query, args...).
		Scan(
			&interview.ID,
			&interview.UserID,
			&interview.QuestionID,
			&interview.StartTimestampMS,
			&interview.EndTimestampMS,
			&interview.Token,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

	args := []any{interview.UserID, interview.QuestionID, interview.StartTimestampMS, interview.EndTimestampMS, interview.Token}

	query := fmt.Sprintf(`
		INSERT INTO %s
			(user_id, question_id, start_timestamp_ms, end_timestamp_ms, token)
		VALUES
			($1, $2, $3, $4, $5)
	`, common.INTERVIEW_TABLE_NAME)

	_, err := i.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to create new interview, %s: %w", err.Error(), common.ErrInternalServerError)
	}

	return nil
}

// Update implements InterviewRepo.
func (i *InterviewRepoImpl) Update(ctx context.Context, interview *entity.Interview) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{interview.UserID, interview.StartTimestampMS, interview.EndTimestampMS, interview.QuestionID, interview.Token, interview.ID}

	query := fmt.Sprintf(`
        UPDATE 
			%s 
        SET 
			user_id = $1,
			start_timestamp_ms = $2,
			end_timestamp_ms = $3,
			question_id = $4,
			token = $5
        WHERE 
			id = $6
    `, common.INTERVIEW_TABLE_NAME)

	result, err := i.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to update interview with id %d, %s: %w", interview.ID, err, common.ErrInternalServerError)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get rows affected when updating interview with id %d, %s: %w", interview.ID, err, common.ErrInternalServerError)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("unable to update interview with id %d: %w", interview.ID, common.ErrInternalServerError)
	}

	return nil
}
