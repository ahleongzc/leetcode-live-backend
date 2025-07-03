package infra

import (
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
)

type InMemoryCallbackQueue interface {
	Enqueue(callback func() error)
	Dequeue() func() error
	Size() uint
}

// The channel has a default size of 100
func NewInMemoryCallbackQueue(
	config *config.InMemoryQueueConfig,
) InMemoryCallbackQueue {
	queueLength := uint(100)

	if config != nil && config.Size != 0 {
		queueLength = config.Size
	}

	return &InMemoryCallbackQueueImpl{
		queue: make(chan func() error, queueLength),
	}
}

type InMemoryCallbackQueueImpl struct {
	queue chan (func() error)
}

func (i *InMemoryCallbackQueueImpl) Enqueue(callback func() error) {
	go func() {
		i.queue <- callback
	}()
}

func (i *InMemoryCallbackQueueImpl) Dequeue() func() error {
	return <-i.queue
}

func (i *InMemoryCallbackQueueImpl) Size() uint {
	return uint(len(i.queue))
}
