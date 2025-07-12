package fasttext

import (
	"context"
	"fmt"
	"sync"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
)

type FastTextPool interface {
	Get(ctx context.Context) (*FastTextProcess, error)
	Put(process *FastTextProcess)
	Classify(ctx context.Context, text string) (*model.IntentDetail, error)
	Close() error
}

type FastTextPoolImpl struct {
	processes chan *FastTextProcess
	modelPath string
	size      uint
	closeOnce sync.Once
}

func NewFastTextPool(
	config *config.IntentClassificationConfig,
) (FastTextPool, error) {
	pool := &FastTextPoolImpl{
		processes: make(chan *FastTextProcess, config.PoolSize),
		modelPath: config.ModelPath,
		size:      config.PoolSize,
	}

	for i := 0; i < int(config.PoolSize); i++ {
		process, err := NewFastTextProcess(config.ModelPath, config.NumClasses)
		if err != nil {
			pool.Close()
			return nil, err
		}
		pool.processes <- process
	}

	return pool, nil
}

func (p *FastTextPoolImpl) Get(ctx context.Context) (*FastTextProcess, error) {
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

func (p *FastTextPoolImpl) Put(process *FastTextProcess) {
	select {
	case p.processes <- process:
	default:
		process.Close()
	}
}

func (p *FastTextPoolImpl) Classify(ctx context.Context, text string) (*model.IntentDetail, error) {
	process, err := p.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer p.Put(process)

	return process.Classify(text)
}

func (p *FastTextPoolImpl) Close() error {
	p.closeOnce.Do(func() {
		close(p.processes)
		for process := range p.processes {
			process.Close()
		}
	})
	return nil
}
