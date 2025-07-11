package repo

import (
	"context"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/fasttext"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type IntentClassificationRepo interface {
	ClassifyIntent(ctx context.Context, word string) (*model.Intent, error)
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

func (i *IntentClassificationRepoImpl) ClassifyIntent(ctx context.Context, word string) (*model.Intent, error) {
	result, err := i.fastTextPool.Classify(ctx, word)
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
