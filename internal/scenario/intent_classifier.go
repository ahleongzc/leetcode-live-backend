package scenario

import (
	"context"
	"strings"
)

type IntentType string

const (
	NO_INTENT             IntentType = "nil"
	HINT_REQUEST          IntentType = "hint_request"
	CLARIFICATION_REQUEST IntentType = "clarification_request"
	END_REQUEST           IntentType = "end_request"
)

type IntentClassifier interface {
	ClassifyIntent(ctx context.Context, message string) (IntentType, error)
}

func NewIntentClassifier() IntentClassifier {
	return &IntentClassifierImpl{}
}

type IntentClassifierImpl struct{}

// TODO: Make an actual intent classifier
func (i *IntentClassifierImpl) ClassifyIntent(ctx context.Context, message string) (IntentType, error) {
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
		"question":      {},
		"questions":     {},
		"confirm":       {},
	}

	// Currently hint takes precedence over clarification
	// TODO: Return the one with higher confidence in the future
	for _, word := range words {
		if _, ok := endKeywords[word]; ok {
			return END_REQUEST, nil
		}

		if _, ok := hintKeywords[word]; ok {
			return HINT_REQUEST, nil
		}

		if _, ok := clarificationRequest[word]; ok {
			return CLARIFICATION_REQUEST, nil
		}
	}

	return NO_INTENT, nil
}
