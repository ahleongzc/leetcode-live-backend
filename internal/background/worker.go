package background

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/infra"

	"github.com/rs/zerolog"
)

type WorkerPool interface {
	Init(size uint)
	Start(ctx context.Context)
}

func NewWorkerPool(
	callbackQueue infra.InMemoryCallbackQueue,
	logger *zerolog.Logger,
) WorkerPool {
	return &WorkerPoolImpl{
		callbackQueue: callbackQueue,
		logger:        logger,
	}
}

type WorkerPoolImpl struct {
	callbackQueue infra.InMemoryCallbackQueue
	size          uint
	logger        *zerolog.Logger
}

func (w *WorkerPoolImpl) Init(size uint) {
	w.size = size
}

func (w *WorkerPoolImpl) Start(ctx context.Context) {
	for workerID := range w.size {
		go func() {
			w.logger.Info().Uint("workerID", workerID).Msg("worker has started")
			for {
				select {
				case <-ctx.Done():
					return
				default:
					w.processCallback(workerID)
				}
			}
		}()
	}
}

func (w *WorkerPoolImpl) processCallback(workerID uint) {
	defer func() {
		if err := recover(); err != nil {
			w.logger.Error().
				Interface("panic", err).
				Uint("workerID", workerID).
				Stack().
				Msg("worker recovered from panic")
		}
	}()

	fn := w.callbackQueue.Dequeue()
	if fn == nil {
		return
	}

	if err := fn(); err != nil {
		w.logger.
			Error().
			Uint("workerID", workerID).
			Err(err).
			Msg("callback error")
	}
}
