package scenario

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra/llm"
	messagequeue "github.com/ahleongzc/leetcode-live-backend/internal/infra/message_queue"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra/tts"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type InterviewScenario interface {
	PrepareToListen(ctx context.Context, interviewID uint) error
	Listen(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error)
	GiveHints(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error)
	Clarify(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error)
	EndInterview(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error)
	ClassifyIntent(ctx context.Context, sentence string) (*model.Intent, error)
}

func NewInterviewScenario(
	reviewScenario ReviewScenario,
	transcriptManager TranscriptManager,
	questionRepo repo.QuestionRepo,
	interviewRepo repo.InterviewRepo,
	fileRepo repo.FileRepo,
	intentClassificationRepo repo.IntentClassificationRepo,
	producer messagequeue.MessageQueueProducer,
	llm llm.LLM,
	tts tts.TTS,
) InterviewScenario {
	return &InterviewScenarioImpl{
		reviewScenario:           reviewScenario,
		transcriptManager:        transcriptManager,
		questionRepo:             questionRepo,
		interviewRepo:            interviewRepo,
		intentClassificationRepo: intentClassificationRepo,
		fileRepo:                 fileRepo,
		producer:                 producer,
		llm:                      llm,
		tts:                      tts,
	}
}

type InterviewScenarioImpl struct {
	reviewScenario           ReviewScenario
	transcriptManager        TranscriptManager
	interviewRepo            repo.InterviewRepo
	questionRepo             repo.QuestionRepo
	fileRepo                 repo.FileRepo
	intentClassificationRepo repo.IntentClassificationRepo
	producer                 messagequeue.MessageQueueProducer
	llm                      llm.LLM
	tts                      tts.TTS
}

// LoadQuestion implements InterviewScenario.
func (i *InterviewScenarioImpl) PrepareToListen(ctx context.Context, interviewID uint) error {
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

func (i *InterviewScenarioImpl) getInterviewQuestionDescription(ctx context.Context, interviewID uint) (string, error) {
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
func (i *InterviewScenarioImpl) Listen(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	return nil, nil
}

// CandidateAsksForClarification implements InterviewScenario.
func (i *InterviewScenarioImpl) Clarify(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	llmMessages := make([]*llm.LLMMessage, 0)
	history, err := i.transcriptManager.GetTranscriptHistory(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	for _, transcript := range history {
		llmMessages = append(llmMessages, transcript.ToLLMMessage())
	}

	llmMessages = append(llmMessages, &llm.LLMMessage{
		Role: llm.SYSTEM,
		Content: `
			You are now tasked to answer clarifying questions from the candidate in a way that helps them better understand the problem without giving away the solution.
			Be clear, concise, and professional — just like you would be in a real interview.
			Provide only as much information as needed to address their question directly.
			Avoid adding extra hints or restating parts of the problem unless it's necessary for clarification.
			If the candidate asks about constraints, edge cases, or assumptions, answer truthfully and succinctly.
			Keep your tone supportive but neutral — you're here to evaluate and guide, not to teach.
		`,
	})

	req := &llm.ChatCompletionsRequest{
		Messages: llmMessages,
	}

	resp, err := i.llm.ChatCompletions(ctx, req)
	if err != nil {
		return nil, err
	}

	replyToCandidate := resp.GetResponse().GetContent()

	reader, err := i.tts.TextToSpeechReader(
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

// CandidateAsksForHints implements InterviewScenario.
func (i *InterviewScenarioImpl) GiveHints(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	llmMessages := make([]*llm.LLMMessage, 0)
	history, err := i.transcriptManager.GetTranscriptHistory(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	llmMessages = append(llmMessages, &llm.LLMMessage{
		Role: llm.SYSTEM,
		Content: `
			You are a senior software engineer conducting a LeetCode-style technical interview. 
			Your task is to provide concise, high-quality hints to help the candidate move forward based on the question they're currently solving and the history of their previous questions or messages. 
			Do not give the full solution. 
			Tailor your hints to their level of understanding and avoid repeating information they've already figured out. 
			If the candidate appears confused or stuck, offer a nudge in the right direction without revealing the answer.
			Keep your hints short and simple, and reply like how you would in a real life interview.
		`,
	})

	for _, transcript := range history {
		llmMessages = append(llmMessages, transcript.ToLLMMessage())
	}

	req := &llm.ChatCompletionsRequest{
		Messages: llmMessages,
	}

	resp, err := i.llm.ChatCompletions(ctx, req)
	if err != nil {
		return nil, err
	}

	replyToCandidate := resp.GetResponse().GetContent()

	reader, err := i.tts.TextToSpeechReader(
		ctx,
		replyToCandidate,
		`You are a senior software engineer conducting a LeetCode-style technical interview. 
		Speak clearly and at a measured pace. Use a calm, thoughtful, and professional tone, as if you're guiding a candidate through the problem. 
		Pause briefly between key points. 
		Avoid sounding robotic—speak naturally and deliberately, like in a real conversation.`,
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

// CandidateWantsToEnd implements InterviewScenario.
func (i *InterviewScenarioImpl) EndInterview(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	interview, err := i.interviewRepo.GetByID(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	interview.EndTimestampMS = util.ToPtr(time.Now().UnixMilli())
	if err = i.interviewRepo.Update(ctx, interview); err != nil {
		return nil, err
	}

	endingTranscript := "Hey thanks for joining this interview, I hope you had fun"
	reader, err := i.tts.TextToSpeechReader(
		ctx,
		endingTranscript,
		`
		You are a senior software engineer conducting a LeetCode-style technical interview. 
		Speak clearly and at a measured pace. Use a calm, thoughtful, and professional tone, as if you're guiding a candidate through the problem. 
		Pause briefly between key points. 
		Avoid sounding robotic—speak naturally and deliberately, like in a real conversation.`,
	)
	if err != nil {
		return nil, err
	}

	reviewMessage, err := json.Marshal(&model.ReviewMessage{InterviewID: interviewID})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal before passing into message queue for review, %s: %w", err, common.ErrInternalServerError)
	}

	if err := i.producer.Push(ctx, reviewMessage, common.REVIEW_QUEUE); err != nil {
		return nil, err
	}

	// TODO: Make the file name follow a structure
	url, err := i.fileRepo.Upload(ctx, "tmp.mp3", reader, nil)
	if err != nil {
		return nil, err
	}

	if err := i.transcriptManager.WriteInterviewer(ctx, interviewID, endingTranscript, url, entity.ASSISTANT); err != nil {
		return nil, err
	}

	return &model.WebSocketMessage{
		From:      model.SERVER,
		URL:       util.ToPtr(url),
		CloseConn: true,
	}, nil
}

// CandidateWantsToEnd implements InterviewScenario.
func (i *InterviewScenarioImpl) AbandonInterview(ctx context.Context, interviewID uint) (*model.WebSocketMessage, error) {
	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	interview, err := i.interviewRepo.GetByID(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	interview.End()

	if err = i.interviewRepo.Update(ctx, interview); err != nil {
		return nil, err
	}

	endingTranscript := "Hey thanks for joining this interview, I hope you had fun"
	reader, err := i.tts.TextToSpeechReader(
		ctx,
		endingTranscript,
		`
		You are a senior software engineer conducting a LeetCode-style technical interview. 
		Speak clearly and at a measured pace. Use a calm, thoughtful, and professional tone, as if you're guiding a candidate through the problem. 
		Pause briefly between key points. 
		Avoid sounding robotic—speak naturally and deliberately, like in a real conversation.`,
	)
	if err != nil {
		return nil, err
	}

	reviewMessage, err := json.Marshal(&model.ReviewMessage{InterviewID: interviewID})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal before passing into message queue for review, %s: %w", err, common.ErrInternalServerError)
	}

	if err := i.producer.Push(ctx, reviewMessage, common.REVIEW_QUEUE); err != nil {
		return nil, err
	}

	// TODO: Make the file name follow a structure
	url, err := i.fileRepo.Upload(ctx, "tmp.mp3", reader, nil)
	if err != nil {
		return nil, err
	}

	if err := i.transcriptManager.WriteInterviewer(ctx, interviewID, endingTranscript, url, entity.ASSISTANT); err != nil {
		return nil, err
	}

	return &model.WebSocketMessage{
		From:      model.SERVER,
		URL:       util.ToPtr(url),
		CloseConn: true,
	}, nil
}

func (i *InterviewScenarioImpl) ClassifyIntent(ctx context.Context, sentence string) (*model.Intent, error) {
	fmt.Println("The sentence is" + sentence)
	intent, err := i.intentClassificationRepo.ClassifyIntent(ctx, sentence)
	if err != nil {
		return nil, err
	}

	return intent, nil
}
