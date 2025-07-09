package repo

import (
	"context"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	intentclassifier "github.com/ahleongzc/leetcode-live-backend/internal/infra/intent_classifier"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type IntentClassificationRepo interface {
	ClassifyIntent(ctx context.Context, word string) (*model.Intent, error)
}

func NewIntentClassificationRepo(
	intentClassifier intentclassifier.IntentClassifier,
) IntentClassificationRepo {
	return &IntentClassificationRepoImpl{
		intentClassifier: intentClassifier,
	}
}

type IntentClassificationRepoImpl struct {
	intentClassifier intentclassifier.IntentClassifier
}

func (i *IntentClassificationRepoImpl) ClassifyIntent(ctx context.Context, word string) (*model.Intent, error) {
	result, err := i.intentClassifier.Classify(ctx, word)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	switch result {
	case "nil", "hint", "end":
		return util.ToPtr(model.Intent(result)), nil
	default:
		return nil, fmt.Errorf("invalid intent %s: %w", result, common.ErrInternalServerError)
	}
}
