package tts

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/go-tts/tts/pkg/speech"
)

type GoTTS struct {
	language string
}

func NewGoTTS(
	language string,
) *GoTTS {
	return &GoTTS{
		language: language,
	}
}

func (g *GoTTS) TextToSpeechReader(ctx context.Context, text, instruction string) (io.Reader, error) {
	// Remove unwanted special characters but keep .,!? and spaces
	re := regexp.MustCompile(`[^a-zA-Z0-9\s.,!?]`)
	text = re.ReplaceAllString(text, "")

	// This is a limitation from the library
	if len(text) > 200 {
		text = text[:200]
	}

	reader, err := speech.FromText(text, g.language)
	if err != nil {
		return nil, fmt.Errorf("unable to generate speech go-tts, %s: %w", err, common.ErrInternalServerError)
	}

	audioData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to read audio data go-tts, %s: %w", err, common.ErrInternalServerError)
	}

	bufferedReader := bytes.NewReader(audioData)
	return bufferedReader, nil
}

func (g *GoTTS) TextToSpeechWriteToFile(ctx context.Context, text, instruction, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("unable to create directory for tts output: %w", common.ErrInternalServerError)
	}

	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create an output file for tts: %w", common.ErrInternalServerError)
	}
	defer outputFile.Close()

	reader, err := speech.FromText(text, g.language)
	if err != nil {
		return fmt.Errorf("unable to generate speech go-tts, %s: %w", err, common.ErrInternalServerError)
	}

	_, err = io.Copy(outputFile, reader)
	if err != nil {
		return fmt.Errorf("unable to write to output file for tts: %w", common.ErrInternalServerError)
	}
	return nil
}
