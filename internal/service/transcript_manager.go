package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type TranscriptManager interface {
	FlushAndRemoveInterview(ctx context.Context, interviewID uint) error
	FlushCandidate(ctx context.Context, interviewID uint) error
	WriteCandidate(ctx context.Context, interviewID uint, chunk string) error
	WriteInterviewer(ctx context.Context, interviewID uint, message, url string) error
	// This sets up the system prompt for the LLM
	PrepareInterviewer(ctx context.Context, interviewID uint, prompt string) error
	GetTranscriptHistory(ctx context.Context, interviewID uint) ([]*entity.Transcript, error)
	GetTranscriptHistoryInLLMMessageFormat(ctx context.Context, interviewID uint) ([]*model.LLMMessage, error)
	HasSufficientWordsInBuffer(ctx context.Context, interviewID uint) (bool, error)
	GetSentenceInBuffer(ctx context.Context, interviewID uint) string
	// Returns the size of the hashmap
	GetManagerInfo() uint
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

// PrepareInterviewer implements TranscriptManager.
func (t *TranscriptManagerImpl) PrepareInterviewer(ctx context.Context, interviewID uint, prompt string) error {
	transcript := entity.NewTranscript().
		SetRole(entity.SYSTEM).
		SetContent(prompt).
		SetInterviewID(interviewID)

	if err := t.transcriptRepo.Create(ctx, transcript); err != nil {
		return err
	}

	return nil
}

func (t *TranscriptManagerImpl) GetTranscriptHistoryInLLMMessageFormat(ctx context.Context, interviewID uint) ([]*model.LLMMessage, error) {
	transcriptHistory, err := t.transcriptRepo.ListByInterviewIDAsc(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	llmMessages := make([]*model.LLMMessage, 0)
	for _, transcript := range transcriptHistory {
		llmMessages = append(llmMessages, transcript.ToLLMMessage())
	}

	return llmMessages, nil
}

// FlushAndRemoveInterview implements TranscriptManager.
func (t *TranscriptManagerImpl) FlushAndRemoveInterview(ctx context.Context, interviewID uint) error {
	if err := t.FlushCandidate(ctx, interviewID); err != nil {
		return err
	}

	delete(t.bufferMap, interviewID)
	return nil
}

// GetManagerInfo implements TranscriptManager.
func (t *TranscriptManagerImpl) GetManagerInfo() uint {
	return uint(len(t.bufferMap))
}

// GetSentenceInBuffer implements TranscriptManager.
func (t *TranscriptManagerImpl) GetSentenceInBuffer(ctx context.Context, interviewID uint) string {
	t.initialiseBuffer(interviewID)
	buffer, ok := t.bufferMap[interviewID]
	if !ok {
		return ""
	}

	return strings.ToLower(strings.TrimSpace(buffer.String()))
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

func (t *TranscriptManagerImpl) WriteInterviewer(ctx context.Context, interviewID uint, chunk, url string) error {
	transcript := entity.NewInterviewerTranscript().
		SetContent(strings.TrimSpace(chunk)).
		SetInterviewID(interviewID).
		SetURL(url)

	err := t.transcriptRepo.Create(ctx, transcript)
	if err != nil {
		return err
	}

	return nil
}

func (t *TranscriptManagerImpl) FlushCandidate(ctx context.Context, interviewID uint) error {
	t.initialiseBuffer(interviewID)
	buffer, ok := t.bufferMap[interviewID]
	if !ok {
		return fmt.Errorf("no buffer exists for interview id %d: %w", interviewID, common.ErrInternalServerError)
	}

	// Nothing to flush into DB
	if len(buffer.String()) == 0 {
		return nil
	}

	trancript := entity.NewCandidateTranscript().
		SetContent(strings.TrimSpace(buffer.String())).
		SetInterviewID(interviewID)

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
