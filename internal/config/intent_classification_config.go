package config

type IntentClassificationConfig struct {
	ModelPath string
	PoolSize  uint
}

func LoadIntentClassificationConfig() (*IntentClassificationConfig, error) {
	modelPath := "./bin/model.bin"
	poolSize := 5

	return &IntentClassificationConfig{
		ModelPath: modelPath,
		PoolSize:  uint(poolSize),
	}, nil
}
