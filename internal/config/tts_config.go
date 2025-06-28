package config

import (
	"os"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

type Voice string
type TTSModel string

const (
	Alloy   Voice = "alloy"
	Ash     Voice = "ash"
	Ballad  Voice = "ballad"
	Coral   Voice = "coral"
	Echo    Voice = "echo"
	Fable   Voice = "fable"
	Nova    Voice = "nova"
	Onyx    Voice = "onyx"
	Sage    Voice = "sage"
	Shimmer Voice = "shimmer"

	TTS_1    TTSModel = "tts-1"
	TTS_1_HD TTSModel = "tts-1-hd"
)

type TTSConfig struct {
	Model  TTSModel
	Voice  Voice
	URL    string
	APIKey string
}

func LoadTTSConfig() *TTSConfig {
	apiKey := os.Getenv(common.OPENAI_API_KEY)

	return &TTSConfig{
		URL:    common.OPENAI_BASE_URL + "/v1/audio/speech",
		Voice:  Alloy,
		Model:  TTS_1,
		APIKey: apiKey,
	}
}
