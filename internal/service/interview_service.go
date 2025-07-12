package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type InterviewService interface {
	GetHistory(ctx context.Context, userID, limit, offset uint) (*model.InterviewHistory, *model.Pagination, error)
	ProcessIncomingMessage(ctx context.Context, interviewID uint, message *model.WebSocketMessage) (*model.WebSocketMessage, error)
	// Returns the id of the interview
	ConsumeTokenAndStartInterview(ctx context.Context, token string) (uint, error)
	// Returns the one-off token that is used to validate the incoming websocket request
	SetUpNewInterview(ctx context.Context, userID uint, externalQuestionID, description string) (string, error)
	PrepareToListen(ctx context.Context, interviewID uint) error
	PauseOngoingInterview(ctx context.Context, interviewID uint) error
	AbandonUnfinishedInterview(ctx context.Context, userID uint) error
	SetUpUnfinishedInterview(ctx context.Context, userID uint) (string, error)
	GetOngoingInterview(ctx context.Context, userID uint) (*model.Interview, error)
	GetUnfinishedInterview(ctx context.Context, userID uint) (*model.Interview, error)
}

func NewInterviewService(
	authService AuthService,
	reviewService ReviewService,
	questionService QuestionService,
	transcriptManager TranscriptManager,
	ttsRepo repo.TTSRepo,
	llmRepo repo.LLMRepo,
	fileRepo repo.FileRepo,
	reviewRepo repo.ReviewRepo,
	questionRepo repo.QuestionRepo,
	interviewRepo repo.InterviewRepo,
	messageQueueRepo repo.MessageQueueProducerRepo,
	intentClassificationRepo repo.IntentClassificationRepo,
) InterviewService {
	return &InterviewServiceImpl{
		authService:              authService,
		reviewService:            reviewService,
		questionService:          questionService,
		transcriptManager:        transcriptManager,
		ttsRepo:                  ttsRepo,
		llmRepo:                  llmRepo,
		fileRepo:                 fileRepo,
		reviewRepo:               reviewRepo,
		questionRepo:             questionRepo,
		interviewRepo:            interviewRepo,
		messageQueueRepo:         messageQueueRepo,
		intentClassificationRepo: intentClassificationRepo,
	}
}

type InterviewServiceImpl struct {
	authService              AuthService
	reviewService            ReviewService
	questionService          QuestionService
	transcriptManager        TranscriptManager
	ttsRepo                  repo.TTSRepo
	llmRepo                  repo.LLMRepo
	fileRepo                 repo.FileRepo
	reviewRepo               repo.ReviewRepo
	questionRepo             repo.QuestionRepo
	interviewRepo            repo.InterviewRepo
	intentClassificationRepo repo.IntentClassificationRepo
	messageQueueRepo         repo.MessageQueueProducerRepo
}

// PauseOngoingInterview implements InterviewService.
func (i *InterviewServiceImpl) PauseOngoingInterview(ctx context.Context, interviewID uint) error {
	interview, err := i.interviewRepo.GetByID(ctx, interviewID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("there is no ongoing interview :%w", common.ErrBadRequest)
		}
		return err
	}

	interview.Pause()
	if err := i.interviewRepo.Update(ctx, interview); err != nil {
		return err
	}

	return nil
}

// GetOngoingInterview implements InterviewService.
func (i *InterviewServiceImpl) GetOngoingInterview(ctx context.Context, userID uint) (*model.Interview, error) {
	ongoingInterview, err := i.interviewRepo.GetOngoingInterviewByUserID(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return nil, err
	}

	if !ongoingInterview.Exists() {
		return nil, nil
	}

	interviewModel, err := i.convertInterviewEntityToModel(ctx, ongoingInterview, nil)
	if err != nil {
		return nil, err
	}

	return interviewModel, nil
}

// AbandonOngoingInterview implements InterviewService.
func (i *InterviewServiceImpl) AbandonUnfinishedInterview(ctx context.Context, userID uint) error {
	interview, err := i.interviewRepo.GetUnfinishedInterviewByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("there is no unfinished interview :%w", common.ErrBadRequest)
		}
		return err
	}

	var errBadRequest error
	if interview.SetUpCount >= 3 {
		return fmt.Errorf("previous set up interview attempt exceeded: %w", common.ErrBadRequest)
	}

	interview.ConsumeToken()
	interview.Abandon()

	if err := i.interviewRepo.Update(ctx, interview); err != nil {
		return err
	}

	if err := i.reviewService.HandleAbandonedInterview(ctx, interview.ID); err != nil {
		return err
	}

	return errBadRequest
}

// SetUpOngoingInterview implements InterviewService.
func (i *InterviewServiceImpl) SetUpUnfinishedInterview(ctx context.Context, userID uint) (string, error) {
	unfinishedInterview, err := i.interviewRepo.GetUnfinishedInterviewByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return "", fmt.Errorf("there is no unfinished interview :%w", common.ErrBadRequest)
		}
		return "", err
	}

	freshToken, err := i.validateInterviewSetUpCount(ctx, unfinishedInterview)
	if err != nil {
		return "", err
	}

	return freshToken, nil
}

func (i *InterviewServiceImpl) GetUnfinishedInterview(ctx context.Context, userID uint) (*model.Interview, error) {
	interview, err := i.interviewRepo.GetUnfinishedInterviewByUserID(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return nil, err
	}

	if !interview.Exists() {
		return nil, nil
	}

	question, err := i.questionRepo.GetByID(ctx, interview.QuestionID)
	if err != nil {
		return nil, err
	}

	interviewModel := &model.Interview{
		ID:                    interview.UUID,
		QuestionAttemptNumber: interview.QuestionAttemptNumber,
		Question:              question.ExternalID,
		StartTimestampS:       util.ToPtr(interview.GetStartTimesampS()),
	}

	return interviewModel, nil
}

func (i *InterviewServiceImpl) getQuestionFromQuestionCacheOrRepo(ctx context.Context, questionID uint, questionCache map[uint]*entity.Question) (*entity.Question, error) {
	if questionCache == nil {
		question, err := i.questionRepo.GetByID(ctx, questionID)
		if err != nil {
			return nil, err
		}
		return question, nil
	}

	if _, ok := questionCache[questionID]; !ok {
		question, err := i.questionRepo.GetByID(ctx, questionID)
		if err != nil {
			return nil, err
		}
		questionCache[questionID] = question
	}

	question, ok := questionCache[questionID]
	if !ok {
		return nil, fmt.Errorf("unable to retrieve question information from cache: %w", common.ErrInternalServerError)
	}

	return question, nil
}

func (i *InterviewServiceImpl) convertInterviewEntityToModel(ctx context.Context, interview *entity.Interview, questionCache map[uint]*entity.Question) (*model.Interview, error) {
	interviewModel := &model.Interview{
		ID:                    interview.UUID,
		QuestionAttemptNumber: interview.QuestionAttemptNumber,
	}

	review, err := i.reviewRepo.GetByID(ctx, interview.GetReviewID())
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return nil, err
	}

	if review != nil {
		interviewModel.Feedback = util.ToPtr(review.Feedback)
		interviewModel.Score = util.ToPtr(review.Score)
		interviewModel.Passed = util.ToPtr(review.Passed)
	}

	question, err := i.getQuestionFromQuestionCacheOrRepo(ctx, interview.QuestionID, questionCache)
	if err != nil {
		return nil, err
	}
	interviewModel.Question = question.ExternalID

	if interview.StartTimestampMS != nil {
		interviewModel.StartTimestampS = util.ToPtr(interview.GetStartTimesampS())
	}

	if interview.EndTimestampMS != nil {
		interviewModel.EndTimestampS = util.ToPtr(interview.GetEndTimestampS())
	}

	return interviewModel, nil
}

// GetHistory implements InterviewService.
func (i *InterviewServiceImpl) GetHistory(ctx context.Context, userID, limit, offset uint) (*model.InterviewHistory, *model.Pagination, error) {
	interviews, total, err := i.interviewRepo.ListStartedInterviewsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	result := make([]*model.Interview, 0)

	questionCache := make(map[uint]*entity.Question)

	for _, interview := range interviews {
		interviewModel, err := i.convertInterviewEntityToModel(ctx, interview, questionCache)
		if err != nil {
			return nil, nil, err
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

func (i *InterviewServiceImpl) SetUpNewInterview(ctx context.Context, userID uint, externalQuestionID, description string) (string, error) {
	questionID, err := i.questionService.GetOrCreateQuestion(ctx, externalQuestionID, description)
	if err != nil {
		return "", err
	}

	unfinishedInterview, err := i.interviewRepo.GetUnfinishedInterviewByUserID(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return "", err
	}

	if unfinishedInterview != nil {
		return "", fmt.Errorf("unfinished interview exists: %w", common.ErrBadRequest)
	}

	unstartedInterview, err := i.interviewRepo.GetUnstartedInterviewByUserID(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return "", err
	}

	if unstartedInterview != nil {
		freshToken, err := i.validateInterviewSetUpCount(ctx, unstartedInterview)
		if err != nil {
			return "", err
		}

		return freshToken, nil
	}

	token := i.authService.GenerateRandomToken()

	questionCount, err := i.interviewRepo.CountByUserIDAndQuestionID(ctx, userID, questionID)
	if err != nil {
		return "", err
	}

	interview := &entity.Interview{
		UserID:                userID,
		QuestionID:            questionID,
		Token:                 util.ToPtr(token),
		QuestionAttemptNumber: questionCount + 1,
		SetUpCount:            1,
	}

	id, err := i.interviewRepo.Create(ctx, interview)
	if err != nil {
		return "", err
	}

	if err := i.PrepareToListen(ctx, id); err != nil {
		return "", nil
	}

	return token, nil
}

func (i *InterviewServiceImpl) ProcessIncomingMessage(ctx context.Context, interviewID uint, message *model.WebSocketMessage) (*model.WebSocketMessage, error) {
	if err := i.transcriptManager.WriteCandidate(ctx, interviewID, util.FromPtr(message.Chunk)); err != nil {
		return nil, err
	}

	sufficient, err := i.transcriptManager.HasSufficientWordsInBuffer(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	if !sufficient {
		return nil, nil
	}

	intent, err := i.intentClassificationRepo.ClassifyIntent(ctx, i.transcriptManager.GetSentenceInBuffer(ctx, interviewID))
	if err != nil {
		return nil, err
	}

	if util.IsDevEnv() {
		intent, score := intent.GetIntentWithHighestConfidenceScoreWithScore()
		fmt.Printf("The current message chunk is '%s', the intent is %s with a score of %f\n", i.transcriptManager.GetSentenceInBuffer(ctx, interviewID), intent, score)
	}

	if err := i.transcriptManager.Flush(ctx, interviewID); err != nil {
		return nil, err
	}

	response, err := i.handleIntent(ctx, interviewID, intent)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *InterviewServiceImpl) handleIntent(ctx context.Context, interviewID uint, intent *model.IntentDetail) (*model.WebSocketMessage, error) {
	if intent == nil {
		return nil, fmt.Errorf("intent cannot be nil: %w", common.ErrInternalServerError)
	}

	if intent.GetIntentWithHighestConfidenceScore() == model.CANDIDATE_EXPLANATION {
		return i.listen(ctx, interviewID)
	}

	if intent.GetIntentWithHighestConfidenceScore() == model.OTHERS {
		return i.answer(ctx, interviewID)
	}

	return nil, fmt.Errorf("invalid intent %v: %w,", util.ToPtr(intent), common.ErrInternalServerError)
}

func (i *InterviewServiceImpl) ConsumeTokenAndStartInterview(ctx context.Context, token string) (uint, error) {
	interview, err := i.interviewRepo.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return 0, common.ErrUnauthorized
		}
		return 0, err
	}

	interview.SetUpCount = 0

	review := &entity.Review{
		Feedback: "The interview is still ongoing...",
	}

	reviewID, err := i.reviewRepo.Create(ctx, review)
	if err != nil {
		return 0, err
	}

	interview.ConsumeToken()

	if !interview.HasStarted() {
		interview.Start()
	}

	interview.SetOngoing()
	interview.ReviewID = util.ToPtr(reviewID)

	err = i.interviewRepo.Update(ctx, interview)
	if err != nil {
		return 0, err
	}

	return interview.ID, nil
}

// Checks if the set up count is more than a certain number then rejects them, else return a fresh token
func (i *InterviewServiceImpl) validateInterviewSetUpCount(ctx context.Context, interview *entity.Interview) (string, error) {
	if interview.SetUpCount >= 3 {
		if interview.Token != nil {
			interview.Token = nil
			if err := i.interviewRepo.Update(ctx, interview); err != nil {
				return "", err
			}
		}

		return "", fmt.Errorf("set up interview attempt exceeded: %w", common.ErrBadRequest)
	}

	token := i.authService.GenerateRandomToken()
	interview.SetUpCount++
	interview.Token = util.ToPtr(token)

	if err := i.interviewRepo.Update(ctx, interview); err != nil {
		return "", err
	}

	return token, nil
}

func (i *InterviewServiceImpl) PrepareToListen(ctx context.Context, interviewID uint) error {
	description, err := i.getInterviewQuestionDescription(ctx, interviewID)
	if err != nil {
		return err
	}

	initialTranscript := fmt.Sprintf(`
		You are a senior software engineer conducting a LeetCode-style technical interview with a candidate.
		You have already prepared the question for the candidate, and the description of the question is as follow:
		%s
	`, description)

	if err := i.transcriptManager.WriteInterviewer(ctx, interviewID, initialTranscript, "", ""); err != nil {
		return err
	}

	return nil
}

func (i *InterviewServiceImpl) getInterviewQuestionDescription(ctx context.Context, interviewID uint) (string, error) {
	interview, err := i.interviewRepo.GetByID(ctx, interviewID)
	if err != nil {
		return "", err
	}

	question, err := i.questionRepo.GetByID(ctx, interview.QuestionID)
	if err != nil {
		return "", err
	}

	return question.Description, nil
}

// ListenToCandidate implements InterviewScenario.
func (i *InterviewServiceImpl) listen(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	return nil, nil
}

// CandidateAsksForClarification implements InterviewScenario.
func (i *InterviewServiceImpl) answer(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	llmMessages := make([]*model.LLMMessage, 0)
	history, err := i.transcriptManager.GetTranscriptHistory(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	for _, transcript := range history {
		llmMessages = append(llmMessages, transcript.ToLLMMessage())
	}

	llmMessages = append(llmMessages, &model.LLMMessage{
		Role: model.SYSTEM,
		Content: `
			You are now tasked to answer any questions from the candidate in a way that helps them better understand the problem without giving away the solution.
			Be clear, concise, and professional — just like you would be in a real interview.
			Provide only as much information as needed to address their question directly.
			Avoid adding extra hints or restating parts of the problem unless it's necessary for clarification.
			If the candidate asks about constraints, edge cases, or assumptions, answer truthfully and succinctly.
			Keep your tone supportive but neutral — you're here to evaluate and guide, not to teach.
		`,
	})

	req := &model.ChatCompletionsRequest{
		Messages: llmMessages,
	}

	resp, err := i.llmRepo.ChatCompletions(ctx, req)
	if err != nil {
		return nil, err
	}

	replyToCandidate := resp.GetResponse().GetContent()

	reader, err := i.ttsRepo.TextToSpeechReader(
		ctx, replyToCandidate, `
			You are a senior software engineer conducting a LeetCode-style technical interview. 
			Speak clearly and at a measured pace. Use a calm, thoughtful, and professional tone, as if you're guiding a candidate through the problem. 
			Pause briefly between key points. 
			Avoid sounding robotic—speak naturally and deliberately, like in a real conversation.
		`,
	)
	if err != nil {
		return nil, err
	}

	// TODO: Make the file name follow a structure
	url, err := i.fileRepo.Upload(ctx, "tmp.mp3", reader, nil)
	if err != nil {
		return nil, err
	}

	if err := i.transcriptManager.WriteInterviewer(ctx, interviewID, replyToCandidate, url, entity.ASSISTANT); err != nil {
		return nil, err
	}

	return &model.WebSocketMessage{
		From: model.SERVER,
		URL:  util.ToPtr(url),
	}, nil
}
