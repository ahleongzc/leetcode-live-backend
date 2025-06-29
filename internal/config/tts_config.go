package config

import (
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

const (
	TTS_DEV_PROVIDER string = "go-tts"
)

type TTSConfig struct {
	Provider string
	Model    string
	Voice    string
	BaseURL  string
	Language string
	APIKey   string
}

func LoadTTSConfig() (*TTSConfig, error) {
	provider := util.GetEnvOr(common.TTS_PROVIDER_KEY, TTS_DEV_PROVIDER)
	model := util.GetEnvOr(common.TTS_MODEL_KEY, "")
	voice := util.GetEnvOr(common.TTS_VOICE_KEY, "")
	baseURL := util.GetEnvOr(common.TTS_BASE_URL_KEY, "")
	apiKey := util.GetEnvOr(common.TTS_API_KEY, "")
	language := util.GetEnvOr(common.TTS_LANGUAGE_KEY, "en")

	if provider != TTS_DEV_PROVIDER {
		if model == "" || voice == "" || baseURL == "" {
			return nil, fmt.Errorf("missing tts config, provider=%s model=%s baseURL=%s: %w", provider, model, baseURL, common.ErrInternalServerError)
		}
		if apiKey == "" {
			return nil, fmt.Errorf("missing api key for provider=%s: %w", provider, common.ErrInternalServerError)
		}
	}

	return &TTSConfig{
		Provider: provider,
		Model:    model,
		Voice:    voice,
		BaseURL:  baseURL,
		APIKey:   apiKey,
		Language: language,
	}, nil
}
