package scenario

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type AuthScenario interface {
	GetUserFromSessionID(ctx context.Context, sessionID string) (*entity.User, error)
}

func NewAuthScenario(
	userRepo repo.UserRepo,
	sessionRepo repo.SessionRepo,
) AuthScenario {
	return &AuthScenarioImpl{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

type AuthScenarioImpl struct {
	userRepo    repo.UserRepo
	sessionRepo repo.SessionRepo
}

func (a *AuthScenarioImpl) GetUserFromSessionID(ctx context.Context, sessionID string) (*entity.User, error) {
	session, err := a.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
