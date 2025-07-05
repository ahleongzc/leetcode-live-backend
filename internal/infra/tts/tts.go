package tts

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
)

type TTS interface {
	TextToSpeechWriteToFile(ctx context.Context, text, instruction, filePath string) error
	TextToSpeechReader(ctx context.Context, text, instruction string) (io.Reader, error)
}

func NewTTS(
	ttsConfig *config.TTSConfig, httpClient *http.Client,
) (TTS, error) {
	switch ttsConfig.Provider {
	case config.TTS_DEV_PROVIDER:
		return NewGoTTS(ttsConfig.Language), nil
	case common.OPENAI:
		return NewOpenAITTS(ttsConfig.Model, ttsConfig.BaseURL, ttsConfig.APIKey, ttsConfig.Voice, httpClient), nil
	default:
		return nil, fmt.Errorf("unsupported TTS provider %s: %w", ttsConfig.Provider, common.ErrInternalServerError)
	}
}
