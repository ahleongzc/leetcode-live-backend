package config

type IntentClassificationConfig struct {
	ModelPath  string
	PoolSize   uint
	NumClasses uint
}

func LoadIntentClassificationConfig() (*IntentClassificationConfig, error) {
	modelPath := "./bin/model.bin"
	numClasses := 2
	poolSize := 5

	return &IntentClassificationConfig{
		ModelPath:  modelPath,
		PoolSize:   uint(poolSize),
		NumClasses: uint(numClasses),
	}, nil
}
