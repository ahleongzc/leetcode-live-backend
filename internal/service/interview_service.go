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
	GetInterviewDetail(ctx context.Context, sessionID, externalQuestionID string) (*entity.Interview, error)
}

func NewInterviewService(
	llm infra.LLM,
	interviewRepo repo.InterviewRepo,
	authScenario scenario.AuthScenario,
) InterviewService {
	return &InterviewServiceImpl{
		llm:           llm,
		interviewRepo: interviewRepo,
		authScenario:  authScenario,
	}
}

type InterviewServiceImpl struct {
	llm           infra.LLM
	interviewRepo repo.InterviewRepo
	authScenario  scenario.AuthScenario
}

func (i *InterviewServiceImpl) ProcessInterviewMessage(ctx context.Context, interviewID int, message *model.InterviewMessage) (*model.InterviewMessage, error) {
	return nil, nil
}

// GetInterviewDetail implements InterviewService.
func (i *InterviewServiceImpl) GetInterviewDetail(ctx context.Context, sessionID, externalQuestionID string) (*entity.Interview, error) {
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

	err = i.interviewRepo.Create(ctx, interview)
	if err != nil {
		return nil, err
	}

	return interview, nil
}
