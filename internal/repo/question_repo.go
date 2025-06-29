package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
)

type QuestionRepo interface {
	Create(ctx context.Context, question *entity.Question) (int, error)
	GetByExternalID(ctx context.Context, externalID string) (*entity.Question, error)
}

func NewQuestionRepo(
	db *sql.DB,
) QuestionRepo {
	return &QuestionRepoImpl{
		db: db,
	}
}

type QuestionRepoImpl struct {
	db *sql.DB
}

// Create implements QuestionRepo.
func (q *QuestionRepoImpl) Create(ctx context.Context, question *entity.Question) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{question.ExternalID, question.Description}

	var id int

	query := fmt.Sprintf(`
		INSERT INTO %s
			(external_id, description)
		VALUES
			($1, $2)
		RETURNING
			id
	`, common.QUESTION_TABLE_NAME)

	err := q.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("unable to create new question, %s: %w", err, common.ErrInternalServerError)
	}

	return id, nil
}

// GetByExternalID implements QuestionRepo.
func (q *QuestionRepoImpl) GetByExternalID(ctx context.Context, externalID string) (*entity.Question, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{externalID}

	query := fmt.Sprintf(`
		SELECT 
			id, external_id, description
		FROM 
			%s
		WHERE 
			external_id = $1
	`, common.QUESTION_TABLE_NAME)

	question := &entity.Question{}
	err := q.db.QueryRowContext(ctx, query, args...).
		Scan(
			&question.ID,
			&question.ExternalID,
			&question.Description,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("question: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get question with external_id %s, %s: %w", externalID, err.Error(), common.ErrInternalServerError)
	}

	return question, nil
}
