package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type TTS interface {
	TextToSpeechWriteToFile(ctx context.Context, text, instruction, filePath string) error
	TextToSpeechReader(ctx context.Context, text, instruction string) (io.Reader, error)
}

func NewTTS(config *config.TTSConfig) TTS {
	return &TTSImpl{
		config:     config,
		httpClient: &http.Client{},
	}
}

type TTSImpl struct {
	config     *config.TTSConfig
	httpClient *http.Client
}

// TextToSpeechGetReader implements TTS.
func (t *TTSImpl) TextToSpeechReader(ctx context.Context, text string, instruction string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(ctx, common.TTS_REQUEST_TIMEOUT)
	defer cancel()

	req, err := t.setUpRequest(ctx, text, instruction)
	if err != nil {
		return nil, err
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform http call for tts: %w", common.ErrInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status for tts is not ok, the status code is %d: %w", resp.StatusCode, common.ErrInternalServerError)
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read audio data, %s: %w", err, common.ErrInternalServerError)
	}

	bufferedReader := bytes.NewReader(audioData)
	return bufferedReader, nil
}

func (t *TTSImpl) TextToSpeechWriteToFile(ctx context.Context, text, instruction, filePath string) error {
	ctx, cancel := context.WithTimeout(ctx, common.TTS_REQUEST_TIMEOUT)
	defer cancel()

	req, err := t.setUpRequest(ctx, text, instruction)
	if err != nil {
		return err
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to perform http call for tts: %w", common.ErrInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status for tts is not ok, the status code is %d: %w", resp.StatusCode, common.ErrInternalServerError)
	}

	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create an output file for tts: %w", common.ErrInternalServerError)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to write to output file for tts: %w", common.ErrInternalServerError)
	}

	return nil
}

func (t *TTSImpl) setUpRequest(ctx context.Context, text, instruction string) (*http.Request, error) {
	payload := util.NewJSONPayload()

	payload.Add("model", t.config.Model)
	payload.Add("input", text)
	payload.Add("voice", t.config.Voice)
	payload.Add("instructions", instruction)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json payload for tts: %w", common.ErrInternalServerError)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.config.URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to generate new request for tts: %w", common.ErrInternalServerError)
	}

	req.Header.Set(common.AUTHORIZATION, fmt.Sprintf("Bearer %s", t.config.APIKey))
	req.Header.Set(common.CONTENT_TYPE, "application/json")
	req.Header.Set(common.ACCEPT, "audio/mpeg")

	return req, nil
}
