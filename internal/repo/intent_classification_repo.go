package repo

import (
	"context"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/fasttext"
)

type IntentClassificationRepo interface {
	// Returns the intent and the confidence score
	ClassifyIntent(ctx context.Context, word string) (*model.IntentDetail, error)
}

func NewIntentClassificationRepo(
	fastTextPool fasttext.FastTextPool,
) IntentClassificationRepo {
	return &IntentClassificationRepoImpl{
		fastTextPool: fastTextPool,
	}
}

type IntentClassificationRepoImpl struct {
	fastTextPool fasttext.FastTextPool
}

func (i *IntentClassificationRepoImpl) ClassifyIntent(ctx context.Context, word string) (*model.IntentDetail, error) {
	intent, err := i.fastTextPool.Classify(ctx, word)
	if err != nil {
		return nil, err
	}

	if _, ok := intent.Mapping[model.CANDIDATE_EXPLANATION]; !ok {
		return nil, fmt.Errorf("missing explanation intent score: %w", common.ErrInternalServerError)
	}

	if _, ok := intent.Mapping[model.OTHERS]; !ok {
		return nil, fmt.Errorf("missing others intent score: %w", common.ErrInternalServerError)
	}

	return intent, nil
}
