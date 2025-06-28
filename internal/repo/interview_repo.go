package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
)

type InterviewQuery struct {
	Where map[string]any
}

type InterviewRepo interface {
	Create(ctx context.Context, interview *entity.Interview) error
	Update(ctx context.Context, interview *entity.Interview) error
	GetByUserIDAndExternalQuestionID(ctx context.Context, userID int, externalQuestionID string) (*entity.Interview, error)
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

func (i *InterviewRepoImpl) GetByUserIDAndExternalQuestionID(ctx context.Context, userID int, externalQuestionID string) (*entity.Interview, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{userID, externalQuestionID}

	query := fmt.Sprintf(`
		SELECT 
			id, user_id, external_question_id, start_timestamp_ms, end_timestamp_ms
		FROM 
			%s
		WHERE 
			user_id = $1 
				AND
			external_question_id = $2
	`, common.INTERVIEW_TABLE_NAME)

	interview := &entity.Interview{}
	err := i.db.QueryRowContext(ctx, query, args...).
		Scan(
			&interview.ID,
			&interview.UserID,
			&interview.ExternalQuestionID,
			&interview.StartTimestampMS,
			&interview.EndTimestampMS,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("interview: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get interview with user id %d and external question id %s, %s: %w", userID, externalQuestionID, err.Error(), common.ErrInternalServerError)
	}

	return interview, nil
}

// Create implements InterviewRepo.
func (i *InterviewRepoImpl) Create(ctx context.Context, interview *entity.Interview) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{interview.UserID, interview.ExternalQuestionID, interview.StartTimestampMS}

	query := fmt.Sprintf(`
		INSERT INTO %s
		(user_id, external_question_id, start_timestamp_ms)
		VALUES
		($1, $2, $3)
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

	args := []any{interview.UserID, interview.ExternalQuestionID, interview.StartTimestampMS, interview.EndTimestampMS, interview.ID}

	query := fmt.Sprintf(`
        UPDATE 
			%s 
        SET 
			user_id = $1,
			external_question_id = $2,
			start_timestamp_ms = $3,
			expire_timestamp_ms = $4
        WHERE 
			id = $5
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
