package tts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type OpenAITTS struct {
	model      string
	baseURL    string
	apiKey     string
	voice      string
	httpClient *http.Client
}

func NewOpenAITTS(
	model, baseURL, apiKey, voice string, httpClient *http.Client,
) *OpenAITTS {
	return &OpenAITTS{
		model:      model,
		baseURL:    baseURL,
		apiKey:     apiKey,
		voice:      voice,
		httpClient: httpClient,
	}
}

func (o *OpenAITTS) TextToSpeechReader(ctx context.Context, text string, instruction string) (io.Reader, error) {
	req, err := o.setUpRequest(ctx, text, instruction)
	if err != nil {
		return nil, err
	}

	resp, err := o.httpClient.Do(req)

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

func (o *OpenAITTS) TextToSpeechWriteToFile(ctx context.Context, text, instruction, filePath string) error {
	req, err := o.setUpRequest(ctx, text, instruction)
	if err != nil {
		return err
	}

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to perform http call for tts: %w", common.ErrInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status for tts is not ok, the status code is %d: %w", resp.StatusCode, common.ErrInternalServerError)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("unable to create directory for tts output: %w", common.ErrInternalServerError)
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

func (o *OpenAITTS) setUpRequest(ctx context.Context, text, instruction string) (*http.Request, error) {
	payload := util.NewJSONPayload()

	payload.Add("model", o.model)
	payload.Add("input", text)
	payload.Add("voice", o.voice)
	payload.Add("instructions", instruction)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json payload for tts: %w", common.ErrInternalServerError)
	}

	url := o.baseURL + "v1/audio/speech"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to generate new request for tts: %w", common.ErrInternalServerError)
	}

	req.Header.Set(common.AUTHORIZATION, fmt.Sprintf("Bearer %s", o.apiKey))
	req.Header.Set(common.CONTENT_TYPE, "application/json")
	req.Header.Set(common.ACCEPT, "audio/mpeg")

	return req, nil
}
