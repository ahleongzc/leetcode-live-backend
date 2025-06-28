package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
)

type TranscriptRepo interface {
	Create(ctx context.Context, transcript *entity.Transcript) error
	ListByInterviewIDAsc(ctx context.Context, interviewID int) ([]*entity.Transcript, error)
	ListByInterviewIDDesc(ctx context.Context, interviewID int) ([]*entity.Transcript, error)
}

func NewTranscriptRepo(
	db *sql.DB,
) TranscriptRepo {
	return &TranscriptRepoImpl{
		db: db,
	}
}

type TranscriptRepoImpl struct {
	db *sql.DB
}

func (t *TranscriptRepoImpl) Create(ctx context.Context, transcript *entity.Transcript) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{transcript.InterviewID, transcript.Role, transcript.Content, transcript.CreatedTimestampMS}

	query := fmt.Sprintf(`
		INSERT INTO %s
		(interview_id, role, content, created_timestamp_ms)
		VALUES
		($1, $2, $3, $4)
	`, common.TRANSCRIPT_TABLE_NAME)

	_, err := t.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to create new transcript, %s: %w", err.Error(), common.ErrInternalServerError)
	}

	return nil
}

func (t *TranscriptRepoImpl) ListByInterviewIDAsc(ctx context.Context, interviewID int) ([]*entity.Transcript, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{interviewID}

	query := fmt.Sprintf(`
        SELECT 
			id, interview_id, role, content, created_timestamp_ms
        FROM 
			%s
        WHERE 
			interview_id = $1
        ORDER BY 
			created_timestamp_ms 
		ASC
    `, common.TRANSCRIPT_TABLE_NAME)

	rows, err := t.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to query for transcript (asc order) for interview id of %d, %s: %w", interviewID, err, common.ErrInternalServerError)
	}
	defer rows.Close()

	var transcripts []*entity.Transcript
	for rows.Next() {
		trancript := &entity.Transcript{}
		if err := rows.Scan(
			&trancript.ID,
			&trancript.InterviewID,
			&trancript.Role,
			&trancript.Content,
			&trancript.CreatedTimestampMS,
		); err != nil {
			return nil, fmt.Errorf("unable to scan transcript (asc order) into struct, %s: %w", err, common.ErrInternalServerError)
		}
		transcripts = append(transcripts, trancript)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("transcript row iteration (asc order) for interview id %d, %s: %w", interviewID, err, common.ErrInternalServerError)
	}

	return transcripts, nil
}

func (t *TranscriptRepoImpl) ListByInterviewIDDesc(ctx context.Context, interviewID int) ([]*entity.Transcript, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{interviewID}

	query := fmt.Sprintf(`
        SELECT 
			id, interview_id, role, content, created_timestamp_ms
        FROM 
			%s
        WHERE 
			interview_id = $1
        ORDER BY 
			created_timestamp_ms 
		DESC
    `, common.TRANSCRIPT_TABLE_NAME)

	rows, err := t.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to query for transcript (desc order) for interview id of %d, %s: %w", interviewID, err, common.ErrInternalServerError)
	}
	defer rows.Close()

	var transcripts []*entity.Transcript
	for rows.Next() {
		trancript := &entity.Transcript{}
		if err := rows.Scan(
			&trancript.ID,
			&trancript.InterviewID,
			&trancript.Role,
			&trancript.Content,
			&trancript.CreatedTimestampMS,
		); err != nil {
			return nil, fmt.Errorf("unable to scan transcript (desc order) into struct, %s: %w", err, common.ErrInternalServerError)
		}
		transcripts = append(transcripts, trancript)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("transcript row iteration (desc order) for interview id %d, %s: %w", interviewID, err, common.ErrInternalServerError)
	}

	return transcripts, nil
}
