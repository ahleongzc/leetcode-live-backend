package intentclassifier

import (
	"context"
	"fmt"
	"sync"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

type FastTextPool struct {
	processes chan *FastTextProcess
	modelPath string
	size      uint
	closeOnce sync.Once
}

func NewFastTextPool(modelPath string, poolSize uint) (*FastTextPool, error) {
	pool := &FastTextPool{
		processes: make(chan *FastTextProcess, poolSize),
		modelPath: modelPath,
		size:      poolSize,
	}

	for i := 0; i < int(poolSize); i++ {
		process, err := NewFastTextProcess(modelPath)
		if err != nil {
			pool.Close()
			return nil, err
		}
		pool.processes <- process
	}

	return pool, nil
}

func (p *FastTextPool) Get(ctx context.Context) (*FastTextProcess, error) {
	select {
	case process, ok := <-p.processes:
		if !ok {
			return nil, fmt.Errorf("pool is closed: %w", common.ErrInternalServerError)
		}
		return process, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", ctx.Err(), common.ErrInternalServerError)
	}
}

func (p *FastTextPool) Put(process *FastTextProcess) {
	select {
	case p.processes <- process:
	default:
		process.Close()
	}
}

func (p *FastTextPool) Classify(ctx context.Context, text string) (string, error) {
	process, err := p.Get(ctx)
	if err != nil {
		return "", err
	}
	defer p.Put(process)

	return process.Classify(text)
}

func (p *FastTextPool) Close() error {
	p.closeOnce.Do(func() {
		close(p.processes)
		for process := range p.processes {
			process.Close()
		}
	})
	return nil
}
