package scenario

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"

	"github.com/google/uuid"
)

type AuthScenario interface {
	GetUserFromSessionToken(ctx context.Context, token string) (*entity.User, error)
	GenerateRandomToken() string
	ValidateSession(ctx context.Context, token string) error
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

// ValidateSession implements AuthScenario.
func (a *AuthScenarioImpl) ValidateSession(ctx context.Context, token string) error {
	if token == "" {
		return common.ErrUnauthorized
	}

	session, err := a.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return common.ErrUnauthorized
		}
		return err
	}

	if session.ExpireTimestampMS < time.Now().UnixMilli() {
		return common.ErrUnauthorized
	}

	session.ExpireTimestampMS = time.Now().Add(48 * time.Hour).UnixMilli()
	err = a.sessionRepo.Update(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthScenarioImpl) GenerateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return uuid.NewString()
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func (a *AuthScenarioImpl) GetUserFromSessionToken(ctx context.Context, token string) (*entity.User, error) {
	err := a.ValidateSession(ctx, token)
	if err != nil {
		return nil, err
	}

	session, err := a.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.ErrUnauthorized
		}
		return nil, err
	}

	user, err := a.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
