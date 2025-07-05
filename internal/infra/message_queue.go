package infra

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	messagequeue "github.com/ahleongzc/leetcode-live-backend/internal/infra/message_queue"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
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
	StartConsuming(ctx context.Context, queue string) (<-chan *model.Delivery, error)
	Close() error
}

func NewMessageQueue(
	config *config.MessageQueueConfig,
) MessageQueue {
	return messagequeue.NewRabbitMQ(config)
}
