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
	WriteCandidate(ctx context.Context, interviewID int, transcript string) error
	WriteInterviewer(ctx context.Context, interviewID int, trascript string) error
	GetTranscriptHistory(ctx context.Context, interviewID int) ([]*entity.Transcript, error)
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

func (t *TranscriptManagerImpl) GetTranscriptHistory(ctx context.Context, interviewID int) ([]*entity.Transcript, error) {
	return nil, nil
}

func (t *TranscriptManagerImpl) WriteInterviewer(ctx context.Context, interviewID int, trascript string) error {
	return nil
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

func (t *TranscriptManagerImpl) initialiseBuffer(ctx context.Context, interviewID int) {
	_, ok := t.bufferMap[interviewID]
	if !ok {
		t.bufferMap[interviewID] = &strings.Builder{}
	}
}

// Write implements TranscriptManager.
func (t *TranscriptManagerImpl) WriteCandidate(ctx context.Context, interviewID int, transcript string) error {
	t.initialiseBuffer(ctx, interviewID)

	buffer, ok := t.bufferMap[interviewID]
	if !ok {
		return fmt.Errorf("buffer is not initialised: %w", common.ErrInternalServerError)
	}

	buffer.WriteString(transcript)

	if buffer.Len() < 128 {
		return nil
	}

	err := t.Flush(ctx, interviewID)
	if err != nil {
		return err
	}

	return nil
}
