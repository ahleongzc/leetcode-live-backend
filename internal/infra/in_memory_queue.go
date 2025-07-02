package infra

import (
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
)

type InMemoryQueue interface {
	Enqueue(message any)
	Dequeue() any
}

// The channel has a default size of 100
func NewInMemoryQueue(
	config *config.InMemoryQueueConfig,
) InMemoryQueue {
	queueLength := 100

	if config.Size != 0 {
		config.Size = uint(queueLength)
	}

	return &InMemoryQueueImpl{
		queue: make(chan any, queueLength),
	}
}

type InMemoryQueueImpl struct {
	queue chan (any)
}

func (i *InMemoryQueueImpl) Enqueue(message any) {
	go func() {
		i.queue <- message
	}()
}

func (i *InMemoryQueueImpl) Dequeue() any {
	return <-i.queue
}
