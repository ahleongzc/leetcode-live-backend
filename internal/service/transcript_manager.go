package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type TranscriptManager interface {
	Flush(ctx context.Context, interviewID uint) error
	WriteCandidate(ctx context.Context, interviewID uint, chunk string) error
	WriteInterviewer(ctx context.Context, interviewID uint, chunk, url string, role entity.Role) error
	GetTranscriptHistory(ctx context.Context, interviewID uint) ([]*entity.Transcript, error)
	HasSufficientWordsInBuffer(ctx context.Context, interviewID uint) (bool, error)
	GetSentenceInBuffer(ctx context.Context, interviewID uint) string
}

func NewTranscriptManager(
	transcriptRepo repo.TranscriptRepo,
) TranscriptManager {
	return &TranscriptManagerImpl{
		transcriptRepo: transcriptRepo,
		bufferMap:      make(map[uint]*strings.Builder),
	}
}

type TranscriptManagerImpl struct {
	transcriptRepo repo.TranscriptRepo
	bufferMap      map[uint]*strings.Builder
}

// GetSentenceInBuffer implements TranscriptManager.
func (t *TranscriptManagerImpl) GetSentenceInBuffer(ctx context.Context, interviewID uint) string {
	t.initialiseBuffer(interviewID)
	buffer, ok := t.bufferMap[interviewID]
	if !ok {
		return ""
	}

	return buffer.String()
}

// WordsInBuffer implements TranscriptManager.
func (t *TranscriptManagerImpl) HasSufficientWordsInBuffer(ctx context.Context, interviewID uint) (bool, error) {
	t.initialiseBuffer(interviewID)
	buffer, ok := t.bufferMap[interviewID]
	if !ok {
		return false, fmt.Errorf("no buffer exists for interview id %d: %w", interviewID, common.ErrInternalServerError)
	}

	if buffer.Len() > 30 {
		return true, nil
	}

	return false, nil
}

func (t *TranscriptManagerImpl) GetTranscriptHistory(ctx context.Context, interviewID uint) ([]*entity.Transcript, error) {
	return t.transcriptRepo.ListByInterviewIDAsc(ctx, interviewID)
}

func (t *TranscriptManagerImpl) WriteInterviewer(ctx context.Context, interviewID uint, chunk, url string, role entity.Role) error {
	trancript := &entity.Transcript{
		Role:        role,
		Content:     strings.TrimSpace(chunk),
		InterviewID: interviewID,
		URL:         url,
	}

	err := t.transcriptRepo.Create(ctx, trancript)
	if err != nil {
		return err
	}

	return nil
}

func (t *TranscriptManagerImpl) Flush(ctx context.Context, interviewID uint) error {
	t.initialiseBuffer(interviewID)
	buffer, ok := t.bufferMap[interviewID]
	if !ok {
		return fmt.Errorf("no buffer exists for interview id %d: %w", interviewID, common.ErrInternalServerError)
	}

	// Nothing to flush into DB
	if len(buffer.String()) == 0 {
		return nil
	}

	trancript := &entity.Transcript{
		Role:        entity.USER,
		Content:     strings.TrimSpace(buffer.String()),
		InterviewID: interviewID,
	}

	err := t.transcriptRepo.Create(ctx, trancript)
	if err != nil {
		return err
	}

	buffer.Reset()
	return nil
}

func (t *TranscriptManagerImpl) initialiseBuffer(interviewID uint) {
	_, ok := t.bufferMap[interviewID]
	if !ok {
		t.bufferMap[interviewID] = &strings.Builder{}
	}
}

// Write implements TranscriptManager.
func (t *TranscriptManagerImpl) WriteCandidate(ctx context.Context, interviewID uint, chunk string) error {
	t.initialiseBuffer(interviewID)
	buffer, ok := t.bufferMap[interviewID]
	if !ok {
		return fmt.Errorf("buffer is not initialised: %w", common.ErrInternalServerError)
	}

	buffer.WriteString(" " + chunk)
	return nil
}
