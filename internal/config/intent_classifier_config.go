package config

type IntentClassifierConfig struct {
	Path string
}

func LoadIntentClassifierConfig() (*IntentClassifierConfig, error) {
	path := "./bin/model.bin"

	return &IntentClassifierConfig{
		Path: path,
	}, nil
}
