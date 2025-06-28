package scenario

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type TranscriptManager interface {
	Flush(ctx context.Context, interviewID int) error
	Write(ctx context.Context, interviewID int, transcript string) error
	InitialiseBuffer(ctx context.Context, interviewID int) error
	DeleteBuffer(ctx context.Context, interviewID int)
}

func NewTranscriptManager(
	transcriptRepo repo.TranscriptRepo,
) TranscriptManager {
	return &TranscriptManagerImpl{
		transcriptRepo: transcriptRepo,
		bufferMap:      make(map[int]*strings.Builder),
	}
}

type TranscriptManagerImpl struct {
	transcriptRepo repo.TranscriptRepo
	bufferMap      map[int]*strings.Builder
}

func (t *TranscriptManagerImpl) DeleteBuffer(ctx context.Context, interviewID int) {
	delete(t.bufferMap, interviewID)
}

func (t *TranscriptManagerImpl) Flush(ctx context.Context, interviewID int) error {
	buffer, ok := t.bufferMap[interviewID]
	if !ok {
		return fmt.Errorf("no buffer exists for interview id %d: %w", interviewID, common.ErrInternalServerError)
	}

	trancript := &entity.Transcript{
		Role:               entity.USER,
		Content:            buffer.String(),
		InterviewID:        interviewID,
		CreatedTimestampMS: time.Now().UnixMilli(),
	}

	err := t.transcriptRepo.Create(ctx, trancript)
	if err != nil {
		return err
	}

	buffer.Reset()
	return nil
}

func (t *TranscriptManagerImpl) InitialiseBuffer(ctx context.Context, interviewID int) error {
	_, ok := t.bufferMap[interviewID]
	if ok {
		return fmt.Errorf("buffer already exists for interview id %d: %w", interviewID, common.ErrInternalServerError)
	}

	t.bufferMap[interviewID] = &strings.Builder{}

	return nil
}

// Write implements TranscriptManager.
func (t *TranscriptManagerImpl) Write(ctx context.Context, interviewID int, transcript string) error {
	buffer, ok := t.bufferMap[interviewID]
	if !ok {
		return fmt.Errorf("no buffer exists for interview id %d: %w", interviewID, common.ErrInternalServerError)
	}

	buffer.WriteString(transcript)

	if buffer.Len() < 1024 {
		return nil
	}

	err := t.Flush(ctx, interviewID)
	if err != nil {
		return err
	}
	return nil
}
