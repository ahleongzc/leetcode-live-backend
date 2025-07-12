package fasttext

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type FastTextProcess struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	scanner *bufio.Scanner
	mu      sync.Mutex
}

func NewFastTextProcess(modelPath string, numClasses uint) (*FastTextProcess, error) {
	cmd := exec.Command("./internal/repo/fasttext/fasttext", "predict-prob", modelPath, "-", strconv.Itoa(int(numClasses)))

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

func (f *FastTextProcess) Classify(text string) (*model.IntentDetail, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if util.ContainsNewline(text) {
		text = strings.ReplaceAll(text, "\n", "")
	}

	if _, err := f.stdin.Write([]byte(text + "\n")); err != nil {
		return nil, fmt.Errorf("failed to write to fasttext: %w", common.ErrInternalServerError)
	}

	if !f.scanner.Scan() {
		if err := f.scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read from fasttext: %w", common.ErrInternalServerError)
		}
		return nil, fmt.Errorf("fasttext process closed unexpectedly: %w", common.ErrInternalServerError)
	}

	parts := strings.Fields(f.scanner.Text())
	if len(parts)%2 != 0 {
		return nil, fmt.Errorf("unexpected fasttext output format %s: %w", f.scanner.Text(), common.ErrInternalServerError)
	}

	intentDetail := model.NewIntentDetail()

	for i := 0; i < len(parts); i += 2 {
		label := strings.TrimPrefix(parts[i], "__label__")

		var confidence float64
		if _, err := fmt.Sscanf(parts[i+1], "%f", &confidence); err != nil {
			return nil, fmt.Errorf("unexpected fasttext output format %s: %w", f.scanner.Text(), common.ErrInternalServerError)
		}
		intentDetail.Mapping[model.Intent(label)] = confidence
	}

	return intentDetail, nil
}

func (f *FastTextProcess) Close() error {
	f.stdin.Close()
	f.stdout.Close()
	return f.cmd.Process.Kill()
}
