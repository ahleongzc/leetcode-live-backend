package messagequeue

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
)

type MessageQueue interface {
	MessageQueueProducer
	MessageQueueConsumer
}

type MessageQueueProducer interface {
	Push(ctx context.Context, data []byte, queue string) error
	Close() error
}

type MessageQueueConsumer interface {
	StartConsuming(ctx context.Context, queue string) (<-chan *Delivery, error)
	Close() error
}

func NewMessageQueue(
	config *config.MessageQueueConfig,
) MessageQueue {
	return NewRabbitMQ(config)
}

type Delivery struct {
	Body []byte
	Acknowledger
}

type Acknowledger interface {
	Ack() error
	Nack(requeue bool) error
	Reject(requeue bool) error
}
