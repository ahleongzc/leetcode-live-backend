package scenario

import (
	"context"
	"strings"

	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type IntentClassifier interface {
	ClassifyIntent(ctx context.Context, message string) (*entity.Intent, error)
}

func NewIntentClassifier() IntentClassifier {
	return &IntentClassifierImpl{}
}

type IntentClassifierImpl struct{}

// TODO: Make an actual intent classifier
func (i *IntentClassifierImpl) ClassifyIntent(ctx context.Context, message string) (*entity.Intent, error) {
	message = strings.ToLower(message)
	words := strings.Split(message, " ")

	endKeywords := map[string]struct{}{
		"bye":     {},
		"bye-bye": {},
	}

	hintKeywords := map[string]struct{}{
		"hint":     {},
		"hints":    {},
		"clue":     {},
		"clues":    {},
		"tip":      {},
		"tips":     {},
		"help":     {},
		"stuck":    {},
		"guidance": {},
	}

	clarificationRequest := map[string]struct{}{
		"ask":           {},
		"clarify":       {},
		"clarifying":    {},
		"clarification": {},
		"confirm":       {},
	}

	// Currently hint takes precedence over clarification
	// TODO: Return the one with higher confidence in the future
	for _, word := range words {
		if _, ok := endKeywords[word]; ok {
			return util.ToPtr(entity.END_REQUEST), nil
		}

		if _, ok := hintKeywords[word]; ok {
			return util.ToPtr(entity.HINT_REQUEST), nil
		}

		if _, ok := clarificationRequest[word]; ok {
			return util.ToPtr(entity.CLARIFICATION_REQUEST), nil
		}
	}

	return util.ToPtr(entity.NO_INTENT), nil
}
