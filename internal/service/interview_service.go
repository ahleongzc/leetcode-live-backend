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
	fileRepo repo.FileRepo,
	authScenario scenario.AuthScenario,
	transcriptManager scenario.TranscriptManager,
	intentClassifier scenario.IntentClassifier,
	tts infra.TTS,
) InterviewService {
	return &InterviewServiceImpl{
		llm:               llm,
		interviewRepo:     interviewRepo,
		authScenario:      authScenario,
		transcriptManager: transcriptManager,
		intentClassifier:  intentClassifier,
		tts:               tts,
		fileRepo:          fileRepo,
	}
}

type InterviewServiceImpl struct {
	llm               infra.LLM
	interviewRepo     repo.InterviewRepo
	fileRepo          repo.FileRepo
	authScenario      scenario.AuthScenario
	transcriptManager scenario.TranscriptManager
	intentClassifier  scenario.IntentClassifier
	tts               infra.TTS
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

	response, err := i.generateReplyBasedOnIntent(ctx, interviewID, intent, errChan)
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

func (i *InterviewServiceImpl) generateReplyBasedOnIntent(ctx context.Context, interviewID int, intent scenario.IntentType, errChan chan error) (*model.InterviewMessage, error) {
	if intent == scenario.NO_INTENT {
		return nil, nil
	}

	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	switch intent {
	case scenario.HINT_REQUEST:
		return i.giveHints(ctx, interviewID, errChan)
	case scenario.CLARIFICATION_REQUEST:
		return i.clarify(ctx, interviewID, errChan)
	case scenario.END_REQUEST:
		return i.endInterview(ctx, interviewID)
	default:
		return nil, nil
	}
}

func (i *InterviewServiceImpl) clarify(ctx context.Context, interviewID int, errChan chan error) (*model.InterviewMessage, error) {
	return nil, nil
}

func (i *InterviewServiceImpl) giveHints(ctx context.Context, interviewID int, errChan chan error) (*model.InterviewMessage, error) {
	history, err := i.transcriptManager.GetTranscriptHistory(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	llmMessages := make([]*model.LLMMessage, len(history)+1)
	llmMessages = append(llmMessages, &model.LLMMessage{
		Role: model.SYSTEM,
		Content: `
		You are a senior software engineer conducting a LeetCode-style technical interview. 
		Your task is to provide concise, high-quality hints to help the candidate move forward based on the question they're currently solving and the history of their previous questions or messages. 
		Do not give the full solution. 
		Tailor your hints to their level of understanding and avoid repeating information they've already figured out. 
		If the candidate appears confused or stuck, offer a nudge in the right direction without revealing the answer.
		Keep your hints short and simple, and reply like how you would in a real life interview.`,
	})

	for _, transcript := range history {
		llmMessages = append(llmMessages, transcript.ToLLMMessage())
	}

	req := &model.ChatCompletionsRequest{
		Messages: llmMessages,
	}

	resp, err := i.llm.ChatCompletions(ctx, req)
	if err != nil {
		return nil, err
	}

	reader, err := i.tts.TextToSpeechReader(
		ctx,
		resp.GetResponse().GetContent(),
		`You are a senior software engineer conducting a LeetCode-style technical interview. 
		Speak clearly and at a measured pace. Use a calm, thoughtful, and professional tone, as if you're guiding a candidate through the problem. 
		Pause briefly between key points. 
		Avoid sounding robotic—speak naturally and deliberately, like in a real conversation.`,
	)
	if err != nil {
		return nil, err
	}

	// TODO: Make the file name follow a structure
	url, err := i.fileRepo.Upload(ctx, "test.mp3", reader, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := i.transcriptManager.WriteInterviewer(ctx, interviewID, resp.GetResponse().GetContent(), url); err != nil {
			errChan <- err
		}
	}()

	return &model.InterviewMessage{
		Type:    model.URL,
		Content: url,
	}, nil
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
