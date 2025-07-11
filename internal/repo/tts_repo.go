package repo

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	gotts "github.com/ahleongzc/leetcode-live-backend/internal/repo/go_tts"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/openai"
)

type TTSRepo interface {
	TextToSpeechWriteToFile(ctx context.Context, text, instruction, filePath string) error
	TextToSpeechReader(ctx context.Context, text, instruction string) (io.Reader, error)
}

func NewTTSRepo(
	ttsConfig *config.TTSConfig, httpClient *http.Client,
) (TTSRepo, error) {
	switch ttsConfig.Provider {
	case config.TTS_DEV_PROVIDER:
		return gotts.NewGoTTS(ttsConfig.Language), nil
	case common.OPENAI:
		return openai.NewOpenAITTS(ttsConfig.Model, ttsConfig.BaseURL, ttsConfig.APIKey, ttsConfig.Voice, httpClient), nil
	default:
		return nil, fmt.Errorf("unsupported TTS provider %s: %w", ttsConfig.Provider, common.ErrInternalServerError)
	}
}
