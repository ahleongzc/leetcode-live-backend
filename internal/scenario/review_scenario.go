package scenario

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra/llm"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type ReviewScenario interface {
	ReviewInterviewPerformance(ctx context.Context, interviewID uint) error
	HandleAbandonedInterview(ctx context.Context, interviewID uint) error
}

func NewReviewScenario(
	reviewRepo repo.ReviewRepo,
	interviewRepo repo.InterviewRepo,
	transcriptManager TranscriptManager,
	llm llm.LLM,
) ReviewScenario {
	return &ReviewScenarioImpl{
		reviewRepo:        reviewRepo,
		interviewRepo:     interviewRepo,
		transcriptManager: transcriptManager,
		llm:               llm,
	}
}

type ReviewScenarioImpl struct {
	reviewRepo        repo.ReviewRepo
	interviewRepo     repo.InterviewRepo
	transcriptManager TranscriptManager
	llm               llm.LLM
}

func (r *ReviewScenarioImpl) HandleAbandonedInterview(ctx context.Context, interviewID uint) error {
	interview, err := r.interviewRepo.GetByID(ctx, interviewID)
	if err != nil {
		return err
	}

	review, err := r.reviewRepo.GetByID(ctx, interview.GetReviewID())
	if err != nil {
		return err
	}

	review.Score = 0
	review.Passed = false
	review.Feedback = "The candidate has abandoned the interview"

	if err := r.updateReviewAndInterview(ctx, review, interview); err != nil {
		return err
	}

	return nil
}

func (r *ReviewScenarioImpl) ReviewInterviewPerformance(ctx context.Context, interviewID uint) error {
	history, err := r.transcriptManager.GetTranscriptHistory(ctx, interviewID)
	if err != nil {
		return err
	}

	llmMessages := make([]*llm.LLMMessage, 0)
	for _, transcript := range history {
		llmMessages = append(llmMessages, transcript.ToLLMMessage())
	}

	llmMessages = append(llmMessages, &llm.LLMMessage{
		Role: llm.ASSISTANT,
		Content: `
			You have now finished conducting the technical interview, and your task is to evaluate the candidate's overall performance based on the history.
			Consider their problem-solving skills, communication clarity, approach to edge cases, code correctness, and ability to respond to hints or clarifications.
			Be objective and professional in your feedback. The feedback is for the candidate to read after they are done.
			Write the feedback as a professional summary intended for the candidate to read after the interview is over — similar to what they would receive in a post-interview review email. 
			It should be concise, formal, and impersonal, without engaging the candidate or inviting further discussion.

			You MUST return a JSON object with the following keys:
			1. 'score': an unsigned integer from 0 to 100 that reflects the candidate's performance.
			2. 'feedback': a concise summary of what the candidate did well and what they could improve.
			3. 'passed': a boolean indicating whether the candidate has passed the interview or not.
		`,
	})

	req := &llm.ChatCompletionsRequest{
		Messages: llmMessages,
	}

	resp, err := r.llm.ChatCompletions(ctx, req)
	if err != nil {
		return err
	}

	llmReviewResponse := &struct {
		Score    uint   `json:"score"`
		Feedback string `json:"feedback"`
		Passed   bool   `json:"passed"`
	}{}

	if err := util.StringToJSON(resp.GetResponse().GetContent(), llmReviewResponse); err != nil {
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

	review.Score = llmReviewResponse.Score
	review.Feedback = llmReviewResponse.Feedback
	review.Passed = llmReviewResponse.Passed

	if err := r.updateReviewAndInterview(ctx, review, interview); err != nil {
		return err
	}

	return nil
}

func (r *ReviewScenarioImpl) updateReviewAndInterview(ctx context.Context, review *entity.Review, interview *entity.Interview) error {
	if err := r.reviewRepo.Update(ctx, review); err != nil {
		return err
	}

	if err := r.interviewRepo.Update(ctx, interview); err != nil {
		return err
	}

	return nil
}
