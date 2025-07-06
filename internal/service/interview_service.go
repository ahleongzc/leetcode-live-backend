package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/scenario"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type InterviewService interface {
	GetHistory(ctx context.Context, userID, limit, offset uint) (*model.InterviewHistory, *model.Pagination, error)
	ProcessIncomingMessage(ctx context.Context, interviewID uint, message *model.WebSocketMessage) (*model.WebSocketMessage, error)
	// Returns the id of the interview
	ConsumeInterviewToken(ctx context.Context, token string) (uint, error)
	// Returns the one-off token that is used to validate the incoming websocket request
	SetUpInterview(ctx context.Context, userID uint, externalQuestionID, description string) (string, error)
}

func NewInterviewService(
	interviewScenario scenario.InterviewScenario,
	authScenario scenario.AuthScenario,
	questionScenario scenario.QuestionScenario,
	intentClassifier scenario.IntentClassifier,
	transcriptManager scenario.TranscriptManager,
	interviewRepo repo.InterviewRepo,
	reviewRepo repo.ReviewRepo,
	questionRepo repo.QuestionRepo,
) InterviewService {
	return &InterviewServiceImpl{
		questionScenario:  questionScenario,
		interviewScenario: interviewScenario,
		authScenario:      authScenario,
		intentClassifier:  intentClassifier,
		transcriptManager: transcriptManager,
		interviewRepo:     interviewRepo,
		reviewRepo:        reviewRepo,
		questionRepo:      questionRepo,
	}
}

type InterviewServiceImpl struct {
	interviewScenario scenario.InterviewScenario
	authScenario      scenario.AuthScenario
	questionScenario  scenario.QuestionScenario
	transcriptManager scenario.TranscriptManager
	intentClassifier  scenario.IntentClassifier
	interviewRepo     repo.InterviewRepo
	reviewRepo        repo.ReviewRepo
	questionRepo      repo.QuestionRepo
}

// GetHistory implements InterviewService.
func (i *InterviewServiceImpl) GetHistory(ctx context.Context, userID, limit, offset uint) (*model.InterviewHistory, *model.Pagination, error) {
	interviews, total, err := i.interviewRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	result := make([]*model.Interview, 0)

	questionCache := make(map[uint]*entity.Question)

	for _, interview := range interviews {
		interviewModel := &model.Interview{
			ID: interview.UUID,
		}

		review, err := i.reviewRepo.GetByID(ctx, interview.GetReviewID())
		if err != nil && !errors.Is(err, common.ErrNotFound) {
			return nil, nil, err
		}

		if review != nil {
			interviewModel.Feedback = util.ToPtr(review.Feedback)
			interviewModel.Score = util.ToPtr(review.Score)
			interviewModel.Passed = util.ToPtr(review.Passed)
		}

		if _, ok := questionCache[interview.QuestionID]; !ok {
			question, err := i.questionRepo.GetByID(ctx, interview.QuestionID)
			if err != nil {
				return nil, nil, err
			}
			questionCache[interview.QuestionID] = question
		}

		question, ok := questionCache[interview.QuestionID]
		if !ok {
			return nil, nil, fmt.Errorf("unable to retrieve question information from cache: %w", common.ErrInternalServerError)
		}

		interviewModel.Question = question.ExternalID

		if interview.StartTimestampMS != nil {
			interviewModel.StartTimestampS = interview.GetStartTimesampS()
		}

		if interview.EndTimestampMS != nil {
			interviewModel.EndTimestampS = util.ToPtr(interview.GetEndTimestampS())
		}

		result = append(result, interviewModel)
	}

	return &model.InterviewHistory{Interviews: result},
		&model.Pagination{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasNext: offset+limit < total,
			HasPrev: offset > 0,
		}, nil
}

func (i *InterviewServiceImpl) SetUpInterview(ctx context.Context, userID uint, externalQuestionID, description string) (string, error) {
	questionID, err := i.questionScenario.GetOrCreateQuestion(ctx, externalQuestionID, description)
	if err != nil {
		return "", err
	}

	ongoingInterview, err := i.interviewScenario.GetOngoingInterview(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return "", err
	}

	if ongoingInterview != nil {
		return "", fmt.Errorf("ongoing interview exists: %w", common.ErrBadRequest)
	}

	token := i.authScenario.GenerateRandomToken()

	interview := &entity.Interview{
		UserID:     userID,
		QuestionID: questionID,
		Token:      util.ToPtr(token),
	}

	id, err := i.interviewRepo.Create(ctx, interview)
	if err != nil {
		return "", err
	}

	if err := i.interviewScenario.PrepareToListen(ctx, id); err != nil {
		return "", nil
	}

	return token, nil
}

func (i *InterviewServiceImpl) ProcessIncomingMessage(ctx context.Context, interviewID uint, message *model.WebSocketMessage) (*model.WebSocketMessage, error) {
	if err := i.transcriptManager.WriteCandidate(ctx, interviewID, util.FromPtr(message.Chunk)); err != nil {
		return nil, err
	}

	intent, err := i.intentClassifier.ClassifyIntent(ctx, util.FromPtr(message.Chunk))
	if err != nil {
		return nil, err
	}

	response, err := i.handleIntent(ctx, interviewID, intent)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *InterviewServiceImpl) handleIntent(ctx context.Context, interviewID uint, intent *entity.Intent) (*model.WebSocketMessage, error) {
	if intent == nil {
		return nil, fmt.Errorf("intent cannot be nil: %w", common.ErrInternalServerError)
	}

	switch util.FromPtr(intent) {
	case entity.NO_INTENT:
		return i.interviewScenario.Listen(ctx, interviewID)
	case entity.HINT_REQUEST:
		return i.interviewScenario.GiveHints(ctx, interviewID)
	case entity.CLARIFICATION_REQUEST:
		return i.interviewScenario.Clarify(ctx, interviewID)
	case entity.END_REQUEST:
		return i.interviewScenario.EndInterview(ctx, interviewID)
	default:
		return nil, fmt.Errorf("invalid intent type %v: %w,", util.ToPtr(intent), common.ErrInternalServerError)
	}
}

func (i *InterviewServiceImpl) ConsumeInterviewToken(ctx context.Context, token string) (uint, error) {
	interview, err := i.interviewRepo.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return 0, common.ErrUnauthorized
		}
		return 0, err
	}

	interview.Token = nil
	interview.StartTimestampMS = util.ToPtr(time.Now().UnixMilli())

	err = i.interviewRepo.Update(ctx, interview)
	if err != nil {
		return 0, err
	}

	return interview.ID, nil
}
