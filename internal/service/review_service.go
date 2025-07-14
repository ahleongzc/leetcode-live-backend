package service

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type ReviewService interface {
	ReviewInterviewPerformance(ctx context.Context, interviewID uint) error
	HandleAbandonedInterview(ctx context.Context, interviewID uint) error
}

func NewReviewService(
	aiUseCase AIUseCase,
	reviewRepo repo.ReviewRepo,
	interviewRepo repo.InterviewRepo,
	transcriptManager TranscriptManager,
) ReviewService {
	return &ReviewServiceImpl{
		aiUseCase:         aiUseCase,
		reviewRepo:        reviewRepo,
		interviewRepo:     interviewRepo,
		transcriptManager: transcriptManager,
	}
}

type ReviewServiceImpl struct {
	aiUseCase         AIUseCase
	reviewRepo        repo.ReviewRepo
	interviewRepo     repo.InterviewRepo
	transcriptManager TranscriptManager
}

func (r *ReviewServiceImpl) HandleAbandonedInterview(ctx context.Context, interviewID uint) error {
	interview, err := r.interviewRepo.GetByID(ctx, interviewID)
	if err != nil {
		return err
	}

	review, err := r.reviewRepo.GetByID(ctx, interview.GetReviewID())
	if err != nil {
		return err
	}

	review.
		SetScore(0).
		SetPassed(false).
		SetFeedback("The candidate has abandoned the interview.")

	if err := r.reviewRepo.Update(ctx, review); err != nil {
		return err
	}

	return nil
}

func (r *ReviewServiceImpl) ReviewInterviewPerformance(ctx context.Context, interviewID uint) error {
	llmMessages, err := r.transcriptManager.GetTranscriptHistoryInLLMMessageFormat(ctx, interviewID)
	if err != nil {
		return err
	}

	prompt := `
		You have now finished conducting the technical interview, and your task is to evaluate the candidate's overall performance based on the history.
		Consider their problem-solving skills, communication clarity, approach to edge cases, code correctness, and ability to respond to hints or clarifications.
		Be objective and professional in your feedback. The feedback is for the candidate to read after they are done.
		Write the feedback as a professional summary intended for the candidate to read after the interview is over — similar to what they would receive in a post-interview review email. 
		It should be concise, formal, and impersonal, without engaging the candidate or inviting further discussion.

		You MUST return a JSON object with the following keys:
		1. 'score': an unsigned integer from 0 to 100 that reflects the candidate's performance.
		2. 'feedback': a concise summary of what the candidate did well and what they could improve.
		3. 'passed': a boolean indicating whether the candidate has passed the interview or not.
	`

	latestPrompt := model.NewLLMMessage().
		SetRole(model.ASSISTANT).
		SetContent(prompt)

	llmMessages = append(llmMessages, latestPrompt)

	reply, err := r.aiUseCase.GenerateTextReply(ctx, llmMessages)
	if err != nil {
		return err
	}

	llmReviewResponse := &struct {
		Score    uint   `json:"score"`
		Feedback string `json:"feedback"`
		Passed   bool   `json:"passed"`
	}{}

	if err := util.StringToJSON(reply, llmReviewResponse); err != nil {
		return err
	}

	interview, err := r.interviewRepo.GetByID(ctx, interviewID)
	if err != nil {
		return err
	}

	review, err := r.reviewRepo.GetByID(ctx, interview.GetReviewID())
	if err != nil {
		return err
	}

	review.
		SetScore(llmReviewResponse.Score).
		SetFeedback(llmReviewResponse.Feedback).
		SetPassed(llmReviewResponse.Passed)

	if err := r.reviewRepo.Update(ctx, review); err != nil {
		return err
	}

	return nil
}
