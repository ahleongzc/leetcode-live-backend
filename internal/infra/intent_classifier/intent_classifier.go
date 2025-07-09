package intentclassifier

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
)

type IntentClassifier interface {
	Classify(ctx context.Context, word string) (string, error)
}

func NewIntentClassifier(
	intentClassifierConfig *config.IntentClassifierConfig,
) (IntentClassifier, error) {
	return NewFastTextPool(intentClassifierConfig.Path, config.INTENT_CLASSIFICATION_MODEL_POOL_SIZE)
}
