package config

import "github.com/ahleongzc/leetcode-live-backend/internal/common"

type InMemoryQueueConfig struct {
	Size uint
}

func LoadInMemoryQueueConfig() (*InMemoryQueueConfig, error) {
	return &InMemoryQueueConfig{
		common.IN_MEMORY_QUEUE_SIZE,
	}, nil
}
