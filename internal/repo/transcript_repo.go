package repo

import (
	"context"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"

	"gorm.io/gorm"
)

type TranscriptRepo interface {
	Create(ctx context.Context, transcript *entity.Transcript) error
	ListByInterviewIDAsc(ctx context.Context, interviewID uint) ([]*entity.Transcript, error)
	ListByInterviewIDDesc(ctx context.Context, interviewID uint) ([]*entity.Transcript, error)
}

func NewTranscriptRepo(
	db *gorm.DB,
) TranscriptRepo {
	return &TranscriptRepoImpl{
		db: db,
	}
}

type TranscriptRepoImpl struct {
	db *gorm.DB
}

func (t *TranscriptRepoImpl) Create(ctx context.Context, transcript *entity.Transcript) error {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	if err := t.db.WithContext(ctx).Create(transcript).Error; err != nil {
		return fmt.Errorf("unable to create new transcript: %w", err)
	}

	return nil
}

func (t *TranscriptRepoImpl) ListByInterviewIDAsc(ctx context.Context, interviewID uint) ([]*entity.Transcript, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	var transcripts []*entity.Transcript
	if err := t.db.WithContext(ctx).
		Where("interview_id = ?", interviewID).
		Order("create_timestamp_ms ASC").
		Find(&transcripts).Error; err != nil {
		return nil, fmt.Errorf("unable to query for transcript (asc order) for interview id of %d, %s: %w", interviewID, err, common.ErrInternalServerError)
	}

	return transcripts, nil
}

func (t *TranscriptRepoImpl) ListByInterviewIDDesc(ctx context.Context, interviewID uint) ([]*entity.Transcript, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DB_QUERY_TIMEOUT)
	defer cancel()

	var transcripts []*entity.Transcript
	if err := t.db.WithContext(ctx).
		Where("interview_id = ?", interviewID).
		Order("create_timestamp_ms DESC").
		Find(&transcripts).Error; err != nil {
		return nil, fmt.Errorf("unable to query for transcript (desc order) for interview id of %d, %s: %w", interviewID, err, common.ErrInternalServerError)
	}

	return transcripts, nil
}
