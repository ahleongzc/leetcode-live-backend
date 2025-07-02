package background

import (
	"sync"

	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/rs/zerolog"
)

type WorkerPool interface {
	Init(size uint)
	Start()
}

func NewReviewWorkerPool(
	queue infra.InMemoryQueue,
	logger *zerolog.Logger,
) WorkerPool {
	return &ReviewWorkerPool{
		queue:  queue,
		logger: logger,
	}
}

type ReviewWorkerPool struct {
	queue  infra.InMemoryQueue
	size   uint
	logger *zerolog.Logger
}

func (r *ReviewWorkerPool) Init(size uint) {
	r.size = size
}

func (r *ReviewWorkerPool) Start() {
	var wg sync.WaitGroup
	for range r.size {
		wg.Add(1)
		go func() {

		}()
	}
}
