package service

import (
	"context"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/scenario"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type InterviewService interface {
	ProcessInterviewMessage(ctx context.Context, interviewID int, message *model.InterviewMessage) (*model.InterviewMessage, error)
	GetInterviewDetails(ctx context.Context, token string) (*entity.Interview, error)
	SetUpInterview(ctx context.Context, sessionID, externalQuestionID, description string) error
}

func NewInterviewService(
	interviewScenario scenario.InterviewScenario,
	authScenario scenario.AuthScenario,
	questionScenario scenario.QuestionScenario,
	transcriptManager scenario.TranscriptManager,
	intentClassifier scenario.IntentClassifier,
	interviewRepo repo.InterviewRepo,
) InterviewService {
	return &InterviewServiceImpl{
		questionScenario:  questionScenario,
		interviewScenario: interviewScenario,
		authScenario:      authScenario,
		transcriptManager: transcriptManager,
		intentClassifier:  intentClassifier,
		interviewRepo:     interviewRepo,
	}
}

type InterviewServiceImpl struct {
	interviewScenario scenario.InterviewScenario
	authScenario      scenario.AuthScenario
	questionScenario  scenario.QuestionScenario
	transcriptManager scenario.TranscriptManager
	intentClassifier  scenario.IntentClassifier
	interviewRepo     repo.InterviewRepo
}

func (i *InterviewServiceImpl) SetUpInterview(ctx context.Context, sessionID, externalQuestionID, description string) error {
	user, err := i.authScenario.GetUserFromSessionID(ctx, sessionID)
	if err != nil {
		return err
	}

	questionID, err := i.questionScenario.GetOrCreateQuestion(ctx, externalQuestionID, description)
	if err != nil {
		return err
	}

	token := i.authScenario.GenerateRandomToken()

	interview := &entity.Interview{
		UserID:           user.ID,
		QuestionID:       questionID,
		StartTimestampMS: time.Now().UnixMilli(),
		Token:            util.ToPtr(token),
	}

	err = i.interviewRepo.Create(ctx, interview)
	if err != nil {
		return err
	}

	return nil
}

func (i *InterviewServiceImpl) ProcessInterviewMessage(ctx context.Context, interviewID int, message *model.InterviewMessage) (*model.InterviewMessage, error) {
	if err := i.transcriptManager.WriteCandidate(ctx, interviewID, message.Content); err != nil {
		return nil, err
	}

	intent, err := i.intentClassifier.ClassifyIntent(ctx, message.Content)
	if err != nil {
		return nil, err
	}

	response, err := i.handleIntent(ctx, interviewID, intent)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *InterviewServiceImpl) handleIntent(ctx context.Context, interviewID int, intent scenario.IntentType) (*model.InterviewMessage, error) {
	if intent == scenario.NO_INTENT {
		return nil, nil
	}

	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	switch intent {
	case scenario.HINT_REQUEST:
		return i.interviewScenario.CandidateAsksForHints(ctx, interviewID)
	case scenario.CLARIFICATION_REQUEST:
		return i.interviewScenario.CandidateAsksForClarification(ctx, interviewID)
	case scenario.END_REQUEST:
		return i.interviewScenario.CandidateWantsToEnd(ctx, interviewID)
	default:
		return nil, nil
	}
}

func (i *InterviewServiceImpl) GetInterviewDetails(ctx context.Context, token string) (*entity.Interview, error) {
	return nil, nil
}
