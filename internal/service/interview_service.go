package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type InterviewService interface {
	ConsumeTokenAndStartInterview(ctx context.Context, token string) (*entity.Interview, error)
	HandleInterviewTimesUp(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error)
	PrepareToListen(ctx context.Context, interviewID uint) error
	PauseCandidateOngoingInterview(ctx context.Context, userID uint) error
	AbandonCandidateUnfinishedInterview(ctx context.Context, userID uint) error
	// Returns the one-off token that is used to validate the incoming websocket request
	SetUpCandidateUnfinishedInterview(ctx context.Context, userID uint) (string, error)
	GetCandidateOngoingInterview(ctx context.Context, userID uint) (*model.Interview, error)
	GetCandidateUnfinishedInterview(ctx context.Context, userID uint) (*model.Interview, error)
	GetHistory(ctx context.Context, userID, limit, offset uint) (*model.InterviewHistory, *model.Pagination, error)
	SetUpNewInterviewForCandidate(ctx context.Context, userID uint, externalQuestionID, description string) (string, error)
	ProcessIncomingMessage(ctx context.Context, interviewID uint, message *model.WebSocketMessage) (*model.WebSocketMessage, error)
	ProcessCandidateMessage(ctx context.Context, interviewID uint, chunk, code string) (*model.InterviewerResponse, error)
}

func NewInterviewService(
	aiUseCase AIUseCase,
	userService UserService,
	authService AuthService,
	reviewService ReviewService,
	questionService QuestionService,
	transcriptManager TranscriptManager,
	fileRepo repo.FileRepo,
	reviewRepo repo.ReviewRepo,
	questionRepo repo.QuestionRepo,
	interviewRepo repo.InterviewRepo,
	messageQueueRepo repo.MessageQueueProducerRepo,
	intentClassificationRepo repo.IntentClassificationRepo,
) InterviewService {
	return &InterviewServiceImpl{
		aiUseCase:                aiUseCase,
		authService:              authService,
		userService:              userService,
		reviewService:            reviewService,
		questionService:          questionService,
		transcriptManager:        transcriptManager,
		fileRepo:                 fileRepo,
		reviewRepo:               reviewRepo,
		questionRepo:             questionRepo,
		interviewRepo:            interviewRepo,
		messageQueueRepo:         messageQueueRepo,
		intentClassificationRepo: intentClassificationRepo,
	}
}

type InterviewServiceImpl struct {
	aiUseCase                AIUseCase
	userService              UserService
	authService              AuthService
	reviewService            ReviewService
	questionService          QuestionService
	transcriptManager        TranscriptManager
	fileRepo                 repo.FileRepo
	reviewRepo               repo.ReviewRepo
	questionRepo             repo.QuestionRepo
	interviewRepo            repo.InterviewRepo
	intentClassificationRepo repo.IntentClassificationRepo
	messageQueueRepo         repo.MessageQueueProducerRepo
}

// ProcessCandidateRequest implements InterviewService.
func (i *InterviewServiceImpl) ProcessCandidateMessage(ctx context.Context, interviewID uint, chunk, code string) (*model.InterviewerResponse, error) {
	if err := i.transcriptManager.WriteCandidate(ctx, interviewID, chunk); err != nil {
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
		intent, score := intent.GetIntentWithHighestConfidenceWithScoreOutOf100()
		if intent == model.OTHERS {
			fmt.Println("!!! Needs to generate reply !!!")
		}
		fmt.Printf("The current message chunk is '%s', the score is %f\n", i.transcriptManager.GetSentenceInBuffer(ctx, interviewID), score)
	}

	if err := i.transcriptManager.FlushCandidate(ctx, interviewID); err != nil {
		return nil, err
	}

	resp, err := i.handleCandidateIntent(ctx, interviewID, intent)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// TimeChecker implements InterviewService.
func (i *InterviewServiceImpl) HandleInterviewTimesUp(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	ongoingInterview, err := i.interviewRepo.GetByID(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	if !ongoingInterview.Exists() {
		return nil, fmt.Errorf("there is no ongoing interview :%w", common.ErrBadRequest)
	}

	ongoingInterview.End()

	if err := i.interviewRepo.Update(ctx, ongoingInterview); err != nil {
		return nil, err
	}

	msg, err := i.timesUp(ctx, interviewID)
	if err != nil {
		return nil, nil
	}

	return msg, nil
}

// PauseOngoingInterview implements InterviewService.
func (i *InterviewServiceImpl) PauseCandidateOngoingInterview(ctx context.Context, userID uint) error {
	interview, err := i.interviewRepo.GetOngoingInterviewByUserID(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return err
	}

	if !interview.Exists() {
		return fmt.Errorf("there is no ongoing interview :%w", common.ErrBadRequest)
	}

	interview.Pause()

	if err := i.transcriptManager.FlushCandidateAndRemoveInterview(ctx, interview.ID); err != nil {
		return err
	}

	if err := i.interviewRepo.Update(ctx, interview); err != nil {
		return err
	}

	return nil
}

func (i *InterviewServiceImpl) GetCandidateOngoingInterview(ctx context.Context, userID uint) (*model.Interview, error) {
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

func (i *InterviewServiceImpl) AbandonCandidateUnfinishedInterview(ctx context.Context, userID uint) error {
	interview, err := i.interviewRepo.GetUnfinishedInterviewByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("there is no unfinished interview :%w", common.ErrBadRequest)
		}
		return err
	}

	var errBadRequest error
	// TODO: Add a check to reset it the next day
	if interview.ExceedSetupCountThreshold() {
		return fmt.Errorf("previous set up interview count exceeded: %w", common.ErrBadRequest)
	}

	interview.Abandon()
	if err := i.interviewRepo.Update(ctx, interview); err != nil {
		return err
	}

	if err := i.transcriptManager.FlushCandidateAndRemoveInterview(ctx, interview.ID); err != nil {
		return err
	}

	if err := i.reviewService.HandleAbandonedInterview(ctx, interview.ID); err != nil {
		return err
	}

	return errBadRequest
}

// SetUpOngoingInterview implements InterviewService.
func (i *InterviewServiceImpl) SetUpCandidateUnfinishedInterview(ctx context.Context, userID uint) (string, error) {
	unfinishedInterview, err := i.interviewRepo.GetUnfinishedInterviewByUserID(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return "", err
	}

	if !unfinishedInterview.Exists() {
		return "", fmt.Errorf("there is no unfinished interview :%w", common.ErrBadRequest)
	}

	freshToken, err := i.validateInterviewSetUpCount(ctx, unfinishedInterview)
	if err != nil {
		return "", err
	}

	return freshToken, nil
}

func (i *InterviewServiceImpl) GetCandidateUnfinishedInterview(ctx context.Context, userID uint) (*model.Interview, error) {
	interview, err := i.interviewRepo.GetUnfinishedInterviewByUserID(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return nil, err
	}

	if !interview.Exists() {
		return nil, nil
	}

	interviewModel, err := i.convertInterviewEntityToModel(ctx, interview, nil)
	if err != nil {
		return nil, err
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
	interviewModel := model.NewInterview().
		SetID(interview.UUID).
		SetQuestionAttemptCount(interview.QuestionAttemptCount).
		SetTimeRemainingS(interview.GetTimeRemainingS())

	review, err := i.reviewRepo.GetByID(ctx, interview.GetReviewID())
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return nil, err
	}

	if review.Exists() {
		interviewModel.
			SetFeedback(review.Feedback).
			SetScore(review.Score).
			SetPassed(review.Passed)
	}

	question, err := i.getQuestionFromQuestionCacheOrRepo(ctx, interview.QuestionID, questionCache)
	if err != nil {
		return nil, err
	}

	interviewModel.SetQuestion(question.ExternalID)

	if interview.HasStarted() {
		interviewModel.SetStartTimestampS(interview.GetStartTimesampS())
	}

	if interview.HasEnded() {
		interviewModel.SetEndTimestampS(interview.GetEndTimestampS())
	}

	return interviewModel, nil
}

// GetHistory implements InterviewService.
func (i *InterviewServiceImpl) GetHistory(ctx context.Context, userID, limit, offset uint) (*model.InterviewHistory, *model.Pagination, error) {
	interviews, total, err := i.interviewRepo.ListStartedInterviewsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	interviewModels := make([]*model.Interview, 0)
	questionCache := make(map[uint]*entity.Question)

	for _, interview := range interviews {
		interviewModel, err := i.convertInterviewEntityToModel(ctx, interview, questionCache)
		if err != nil {
			return nil, nil, err
		}
		interviewModels = append(interviewModels, interviewModel)
	}

	pagination := model.NewPagination().
		SetTotal(total).
		SetLimit(limit).
		SetOffset(offset).
		SetHasNext(offset+limit < total).
		SetHasPrev(offset > 0)

	history := model.NewInterviewHistory().
		SetInterviews(interviewModels)

	return history, pagination, nil
}

func (i *InterviewServiceImpl) SetUpNewInterviewForCandidate(ctx context.Context, userID uint, externalQuestionID, description string) (string, error) {
	questionID, err := i.questionService.GetOrCreateQuestion(ctx, externalQuestionID, description)
	if err != nil {
		return "", err
	}

	unfinishedInterview, err := i.interviewRepo.GetUnfinishedInterviewByUserID(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return "", err
	}

	if unfinishedInterview.Exists() {
		return "", fmt.Errorf("unfinished interview exists: %w", common.ErrBadRequest)
	}

	unstartedInterview, err := i.interviewRepo.GetUnstartedInterviewByUserID(ctx, userID)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return "", err
	}

	if unstartedInterview.Exists() {
		freshToken, err := i.validateInterviewSetUpCount(ctx, unstartedInterview)
		if err != nil {
			return "", err
		}

		return freshToken, nil
	}

	questionCount, err := i.interviewRepo.CountByUserIDAndQuestionID(ctx, userID, questionID)
	if err != nil {
		return "", err
	}

	setting, err := i.userService.GetUserSetting(ctx, userID)
	if err != nil {
		return "", err
	}

	interview := entity.NewInterview().
		SetUserID(userID).
		SetQuestionID(questionID).
		SetToken(i.authService.GenerateRandomToken()).
		SetQuestionAttemptCount(questionCount + 1).
		SetSetupCount(1).
		SetAllocatedDurationS(setting.InterviewDurationS)

	id, err := i.interviewRepo.Create(ctx, interview)
	if err != nil {
		return "", err
	}

	if err := i.PrepareToListen(ctx, id); err != nil {
		return "", nil
	}

	return interview.GetToken(), nil
}

// TODO: Add a new method here to process message that are in the buffer after certain delay for better user experience
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
		intent, score := intent.GetIntentWithHighestConfidenceWithScoreOutOf100()
		if intent == model.OTHERS {
			fmt.Println("!!! Needs to generate reply !!!")
		}
		fmt.Printf("The current message chunk is '%s', the score is %f\n", i.transcriptManager.GetSentenceInBuffer(ctx, interviewID), score)
	}

	if err := i.transcriptManager.FlushCandidate(ctx, interviewID); err != nil {
		return nil, err
	}

	response, err := i.handleIntent(ctx, interviewID, intent)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *InterviewServiceImpl) handleCandidateIntent(ctx context.Context, interviewID uint, intentDetail *model.IntentDetail) (*model.InterviewerResponse, error) {
	if !intentDetail.Exists() {
		return nil, fmt.Errorf("intent cannot be nil: %w", common.ErrInternalServerError)
	}

	intent, score := intentDetail.GetIntentWithHighestConfidenceWithScoreOutOf100()
	if intent == model.CANDIDATE_EXPLANATION {
		return i.listenToCandidate(ctx, interviewID)
	}

	// Others mean you would need to answer back, the candidate might be asking for clarification or hints etc etc
	// This will only be triggered if the score is more than 70 out of 100
	if intent == model.OTHERS {
		if score > 70 {
			return i.answerCandidate(ctx, interviewID)
		}
		return i.listenToCandidate(ctx, interviewID)
	}

	return nil, fmt.Errorf("invalid intent %v: %w,", util.ToPtr(intent), common.ErrInternalServerError)
}

func (i *InterviewServiceImpl) handleIntent(ctx context.Context, interviewID uint, intentDetail *model.IntentDetail) (*model.WebSocketMessage, error) {
	if !intentDetail.Exists() {
		return nil, fmt.Errorf("intent cannot be nil: %w", common.ErrInternalServerError)
	}

	intent, score := intentDetail.GetIntentWithHighestConfidenceWithScoreOutOf100()
	if intent == model.CANDIDATE_EXPLANATION {
		return i.listen(ctx, interviewID)
	}

	// Others mean you would need to answer back, the candidate might be asking for clarification or hints etc etc
	// This will only be triggered if the score is more than 70 / 100
	if intent == model.OTHERS {
		if score > 70 {
			return i.answer(ctx, interviewID)
		}
		return i.listen(ctx, interviewID)
	}

	return nil, fmt.Errorf("invalid intent %v: %w,", util.ToPtr(intent), common.ErrInternalServerError)
}

// The validation have to happen here as websockets can only pass the token in query params
func (i *InterviewServiceImpl) ConsumeTokenAndStartInterview(ctx context.Context, token string) (*entity.Interview, error) {
	interview, err := i.interviewRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if !interview.Exists() {
		return nil, common.ErrUnauthorized
	}

	if !interview.ReviewExists() {
		review := entity.NewReview().
			SetFeedback("The interview is still ongoing")

		reviewID, err := i.reviewRepo.Create(ctx, review)
		if err != nil {
			return nil, err
		}

		interview.SetReviewID(reviewID)
	}

	if !interview.HasStarted() {
		interview.Start()
	}

	interview.
		ConsumeToken().
		SetOngoing().
		ResetSetupCount()

	err = i.interviewRepo.Update(ctx, interview)
	if err != nil {
		return nil, err
	}

	return interview, nil
}

// Checks if the set up count is more than a certain number then rejects them, else return a fresh token
func (i *InterviewServiceImpl) validateInterviewSetUpCount(ctx context.Context, interview *entity.Interview) (string, error) {
	if interview.ExceedSetupCountThreshold() {
		if interview.TokenExists() {
			interview.ConsumeToken()
			if err := i.interviewRepo.Update(ctx, interview); err != nil {
				return "", err
			}
		}

		return "", fmt.Errorf("set up interview attempt exceeded: %w", common.ErrBadRequest)
	}

	token := i.authService.GenerateRandomToken()

	interview.
		IncrementSetupCount().
		SetToken(token)

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

	prompt := fmt.Sprintf(`
		You are a senior software engineer conducting a LeetCode-style technical interview with a candidate.
		You have already prepared the question for the candidate, and the description of the question is as follow:
		%s
	`, description)

	if err := i.transcriptManager.PrepareInterviewer(ctx, interviewID, prompt); err != nil {
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

// TODO: Add more functionalities here like logging?
func (i *InterviewServiceImpl) listenToCandidate(ctx context.Context, interviewID uint) (*model.InterviewerResponse, error) {
	return nil, nil
}

// TODO: Add more functionalities here like logging?
func (i *InterviewServiceImpl) listen(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	return nil, nil
}

func (i *InterviewServiceImpl) answerCandidate(ctx context.Context, interviewID uint) (*model.InterviewerResponse, error) {
	prompt := `
		Based on the transcript history, you are now tasked to reply to the candidate.
		If the candidate is asking for a hint, you must answer in a way that helps them better understand the problem without giving away the solution.
		Provide only as much information as needed to address their question directly.
		Avoid adding extra hints or restating parts of the problem unless it's necessary for clarification.
		If the candidate asks about constraints, clarification on the problem, edge cases, or assumptions, answer truthfully and succinctly.
		Be clear, concise, and professional — just like you would be in a real interview.
	`
	replyToCandidate, err := i.generateTextReply(ctx, prompt, interviewID)
	if err != nil {
		return nil, err
	}

	reader, err := i.generateSpeechReply(ctx, replyToCandidate)
	if err != nil {
		return nil, err
	}

	url, err := i.uploadVoiceReply(ctx, interviewID, reader)
	if err != nil {
		return nil, err
	}

	if err := i.transcriptManager.WriteInterviewer(ctx, interviewID, replyToCandidate, url); err != nil {
		return nil, err
	}

	resp := model.NewInterviewerResponse().
		SetURL(url)

	return resp, nil
}

func (i *InterviewServiceImpl) timesUp(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	if err := i.transcriptManager.FlushCandidateAndRemoveInterview(ctx, interviewID); err != nil {
		return nil, err
	}

	prompt := `
		The interview time is up.
		You need to thank the candidate for their time for attempting the interview and remind them that you will review the entire process and give him an appropriate score later.
		Be clear, concise, and professional — just like you would be in a real interview.
	`

	replyToCandidate, err := i.generateTextReply(ctx, prompt, interviewID)
	if err != nil {
		return nil, err
	}

	reader, err := i.generateSpeechReply(ctx, replyToCandidate)
	if err != nil {
		return nil, err
	}

	url, err := i.uploadVoiceReply(ctx, interviewID, reader)
	if err != nil {
		return nil, err
	}

	if err := i.transcriptManager.WriteInterviewer(ctx, interviewID, replyToCandidate, url); err != nil {
		return nil, err
	}

	msg := model.NewServerWebsocketMessage().
		SetURL(url).
		CloseConnection()

	return msg, nil
}

// CandidateAsksForClarification implements InterviewScenario.
func (i *InterviewServiceImpl) answer(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	prompt := `
		Based on the transcript history, you are now tasked to reply to the candidate.
		If the candidate is asking for a hint, you must answer in a way that helps them better understand the problem without giving away the solution.
		Provide only as much information as needed to address their question directly.
		Avoid adding extra hints or restating parts of the problem unless it's necessary for clarification.
		If the candidate asks about constraints, clarification on the problem, edge cases, or assumptions, answer truthfully and succinctly.
		Be clear, concise, and professional — just like you would be in a real interview.
	`
	replyToCandidate, err := i.generateTextReply(ctx, prompt, interviewID)
	if err != nil {
		return nil, err
	}

	reader, err := i.generateSpeechReply(ctx, replyToCandidate)
	if err != nil {
		return nil, err
	}

	url, err := i.uploadVoiceReply(ctx, interviewID, reader)
	if err != nil {
		return nil, err
	}

	if err := i.transcriptManager.WriteInterviewer(ctx, interviewID, replyToCandidate, url); err != nil {
		return nil, err
	}

	msg := model.NewServerWebsocketMessage().
		SetURL(url)

	return msg, nil
}

func (i *InterviewServiceImpl) generateSpeechReply(ctx context.Context, content string) (io.Reader, error) {
	instruction := `
		You are a senior software engineer conducting a LeetCode-style technical interview. 
		Speak clearly and at a measured pace. Use a calm, thoughtful, and professional tone, as if you're guiding a candidate through the problem. 
		Pause briefly between key points. 
		Avoid sounding robotic—speak naturally and deliberately, like in a real conversation.
	`

	reader, err := i.aiUseCase.GenerateSpeechReply(ctx, content, instruction)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (i *InterviewServiceImpl) uploadVoiceReply(ctx context.Context, interviewID uint, reader io.Reader) (string, error) {
	defer func() {
		if util.IsDevEnv() {
			fmt.Println("Finished sending reply to frontend")
		}
	}()

	interview, err := i.interviewRepo.GetByID(ctx, interviewID)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("user_%d/interview_%d/timestamp_ms_%d.mp3", interview.UserID, interviewID, time.Now().UnixMilli())

	url, err := i.fileRepo.Upload(ctx, path, reader, nil)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (i *InterviewServiceImpl) generateTextReply(ctx context.Context, prompt string, interviewID uint) (string, error) {
	llmMessages, err := i.transcriptManager.GetTranscriptHistoryInLLMMessageFormat(ctx, interviewID)
	if err != nil {
		return "", err
	}

	latestPrompt := model.NewLLMMessage().
		SetRole(model.ASSISTANT).
		SetContent(prompt)

	llmMessages = append(llmMessages, latestPrompt)

	replyToCandidate, err := i.aiUseCase.GenerateTextReply(ctx, llmMessages)
	if err != nil {
		return "", err
	}

	return replyToCandidate, nil
}
