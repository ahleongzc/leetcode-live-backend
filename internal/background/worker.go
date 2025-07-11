package background

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/repo"

	"github.com/rs/zerolog"
)

type WorkerPool interface {
	Init(size uint)
	Start(ctx context.Context)
}

func NewWorkerPool(
	callbackQueueRepo repo.InMemoryCallbackQueueRepo,
	logger *zerolog.Logger,
) WorkerPool {
	return &WorkerPoolImpl{
		callbackQueueRepo: callbackQueueRepo,
		logger:            logger,
	}
}

type WorkerPoolImpl struct {
	callbackQueueRepo repo.InMemoryCallbackQueueRepo
	size              uint
	logger            *zerolog.Logger
}

func (w *WorkerPoolImpl) Init(size uint) {
	w.size = size
}

func (w *WorkerPoolImpl) Start(ctx context.Context) {
	for workerID := range w.size {
		go func() {
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

	fn := w.callbackQueueRepo.Dequeue()
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

	w.logger.Info().Uint("workerID", workerID).Msg("worker has finished executing callback function")
}
