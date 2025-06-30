package scenario

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type InterviewScenario interface {
	Listen(ctx context.Context, interviewID int) (*model.InterviewMessage, error)
	GiveHints(ctx context.Context, interviewID int) (*model.InterviewMessage, error)
	Clarify(ctx context.Context, interviewID int) (*model.InterviewMessage, error)
	EndInterview(ctx context.Context, interviewID int) (*model.InterviewMessage, error)
}

func NewInterviewScenario(
	transcriptManager TranscriptManager,
	questionRepo repo.QuestionRepo,
	fileRepo repo.FileRepo,
	llm infra.LLM,
	tts infra.TTS,
) InterviewScenario {
	return &InterviewScenarioImpl{
		transcriptManager: transcriptManager,
		questionRepo:      questionRepo,
		fileRepo:          fileRepo,
		llm:               llm,
		tts:               tts,
	}
}

type InterviewScenarioImpl struct {
	transcriptManager TranscriptManager
	questionRepo      repo.QuestionRepo
	fileRepo          repo.FileRepo
	llm               infra.LLM
	tts               infra.TTS
}

// ListenToCandidate implements InterviewScenario.
func (i *InterviewScenarioImpl) Listen(ctx context.Context, interviewID int) (*model.InterviewMessage, error) {
	return nil, nil
}

// CandidateAsksForClarification implements InterviewScenario.
func (i *InterviewScenarioImpl) Clarify(ctx context.Context, interviewID int) (*model.InterviewMessage, error) {
	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// CandidateAsksForHints implements InterviewScenario.
func (i *InterviewScenarioImpl) GiveHints(ctx context.Context, interviewID int) (*model.InterviewMessage, error) {
	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}

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
		Keep your hints short and simple, and reply like how you would in a real life interview.
		The question that the candidate is solving is Two Sum`,
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

	if err := i.transcriptManager.WriteInterviewer(ctx, interviewID, replyToCandidate, url); err != nil {
		return nil, err
	}

	return &model.InterviewMessage{
		Type:    model.URL,
		Content: url,
	}, nil
}

// CandidateWantsToEnd implements InterviewScenario.
func (i *InterviewScenarioImpl) EndInterview(ctx context.Context, interviewID int) (*model.InterviewMessage, error) {
	err := i.transcriptManager.Flush(ctx, interviewID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
