package service

import (
	"context"
	"errors"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/scenario"
)

type InterviewService interface {
	ProcessInterviewMessage(ctx context.Context, interviewID int, message *model.InterviewMessage) (*model.InterviewMessage, error)
	GetInterviewDetails(ctx context.Context, sessionID, externalQuestionID string) (*entity.Interview, error)
}

func NewInterviewService(
	llm infra.LLM,
	interviewRepo repo.InterviewRepo,
	authScenario scenario.AuthScenario,
	transcriptManager scenario.TranscriptManager,
	intentClassifier scenario.IntentClassifier,
) InterviewService {
	return &InterviewServiceImpl{
		llm:               llm,
		interviewRepo:     interviewRepo,
		authScenario:      authScenario,
		transcriptManager: transcriptManager,
		intentClassifier:  intentClassifier,
	}
}

type InterviewServiceImpl struct {
	llm               infra.LLM
	interviewRepo     repo.InterviewRepo
	authScenario      scenario.AuthScenario
	transcriptManager scenario.TranscriptManager
	intentClassifier  scenario.IntentClassifier
}

func (i *InterviewServiceImpl) ProcessInterviewMessage(ctx context.Context, interviewID int, message *model.InterviewMessage) (*model.InterviewMessage, error) {
	errChan := make(chan error, 1)
	go func() {
		if err := i.transcriptManager.WriteCandidate(ctx, interviewID, message.Content); err != nil {
			errChan <- err
		}
	}()

	intent, err := i.intentClassifier.ClassifyIntent(ctx, message.Content)
	if err != nil {
		return nil, err
	}

	response, err := i.generateReplyBasedOnIntent(ctx, interviewID, intent)
	if err != nil {
		return nil, err
	}

	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	default:
	}

	return response, nil
}

func (i *InterviewServiceImpl) generateReplyBasedOnIntent(ctx context.Context, interviewID int, intent scenario.IntentType) (*model.InterviewMessage, error) {
	switch intent {
	case scenario.NO_INTENT:
		return nil, nil
	case scenario.HINT_REQUEST:
		return i.giveHints(ctx, interviewID)
	case scenario.CLARIFICATION_REQUEST:
		return i.clarify(ctx, interviewID)
	case scenario.END_REQUEST:
		return i.endInterview(ctx, interviewID)
	default:
		return nil, nil
	}
}

func (i *InterviewServiceImpl) clarify(ctx context.Context, interviewID int) (*model.InterviewMessage, error) {
	return nil, nil
}

func (i *InterviewServiceImpl) giveHints(ctx context.Context, interviewID int) (*model.InterviewMessage, error) {
	return nil, nil
}

func (i *InterviewServiceImpl) endInterview(ctx context.Context, interviewID int) (*model.InterviewMessage, error) {
	return nil, nil
}

func (i *InterviewServiceImpl) GetInterviewDetails(ctx context.Context, sessionID, externalQuestionID string) (*entity.Interview, error) {
	user, err := i.authScenario.GetUserFromSessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	interview, err := i.interviewRepo.GetByUserIDAndExternalQuestionID(ctx, user.ID, externalQuestionID)
	if nil == err && interview != nil {
		return interview, nil
	}

	if err != nil {
		if !errors.Is(err, common.ErrNotFound) {
			return nil, err
		}
	}

	interview = &entity.Interview{
		UserID:             user.ID,
		ExternalQuestionID: externalQuestionID,
		StartTimestampMS:   time.Now().UnixMilli(),
	}

	if err := i.interviewRepo.Create(ctx, interview); err != nil {
		return nil, err
	}

	return interview, nil
}
