package repo

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/rabbitmq.go"
)

type MessageQueueRepo interface {
	MessageQueueProducerRepo
	MessageQueueConsumerRepo
}

type MessageQueueProducerRepo interface {
	Push(ctx context.Context, data []byte, queue string) error
	Close() error
}

type MessageQueueConsumerRepo interface {
	StartConsuming(ctx context.Context, queue string) (<-chan *model.Delivery, error)
	Close() error
}

func NewMessageQueueRepo(
	config *config.MessageQueueConfig,
) MessageQueueRepo {
	return rabbitmq.NewRabbitMQ(config)
}
