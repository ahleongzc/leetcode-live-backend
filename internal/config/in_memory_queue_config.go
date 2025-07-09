package config

type InMemoryQueueConfig struct {
	Size uint
}

func LoadInMemoryQueueConfig() (*InMemoryQueueConfig, error) {
	return &InMemoryQueueConfig{
		IN_MEMORY_QUEUE_SIZE,
	}, nil
}
