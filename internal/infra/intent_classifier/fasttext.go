package intentclassifier

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type FastTextProcess struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	scanner *bufio.Scanner
	mu      sync.Mutex
}

func NewFastTextProcess(modelPath string) (*FastTextProcess, error) {
	cmd := exec.Command("./internal/infra/intent_classifier/fasttext", "predict", modelPath, "-")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", common.ErrInternalServerError)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", common.ErrInternalServerError)
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		return nil, fmt.Errorf("failed to start fasttext process, %s: %w", err, common.ErrInternalServerError)
	}

	return &FastTextProcess{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		scanner: bufio.NewScanner(stdout),
	}, nil
}

func (f *FastTextProcess) Classify(text string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if util.ContainsNewline(text) {
		text = strings.ReplaceAll(text, "\n", "")
	}

	if _, err := f.stdin.Write([]byte(text + "\n")); err != nil {
		return "", fmt.Errorf("failed to write to fasttext: %w", common.ErrInternalServerError)
	}

	if !f.scanner.Scan() {
		if err := f.scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read from fasttext: %w", common.ErrInternalServerError)
		}
		return "", fmt.Errorf("fasttext process closed unexpectedly: %w", common.ErrInternalServerError)
	}

	return strings.TrimPrefix(f.scanner.Text(), "__label__"), nil
}

func (f *FastTextProcess) Close() error {
	f.stdin.Close()
	f.stdout.Close()
	return f.cmd.Process.Kill()
}
