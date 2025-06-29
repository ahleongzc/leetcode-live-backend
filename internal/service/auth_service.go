package service

import (
	"context"
	"errors"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/scenario"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, sessionID string) error
}

func NewAuthService(
	authScenario scenario.AuthScenario,
	sessionRepo repo.SessionRepo,
	userRepo repo.UserRepo,
) AuthService {
	return &AuthServiceImpl{
		authScenario: authScenario,
		sessionRepo:  sessionRepo,
		userRepo:     userRepo,
	}
}

type AuthServiceImpl struct {
	authScenario scenario.AuthScenario
	sessionRepo  repo.SessionRepo
	userRepo     repo.UserRepo
}

func (a *AuthServiceImpl) Logout(ctx context.Context, sessionID string) error {
	err := a.sessionRepo.DeleteByID(ctx, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return "", common.ErrUnauthorized
		}
		return "", err
	}

	if !isSamePassword(user.Password, password) {
		return "", common.ErrUnauthorized
	}

	session := &entity.Session{
		ID:                a.authScenario.GenerateRandomToken(),
		UserID:            user.ID,
		ExpireTimestampMS: time.Now().Add(48 * time.Hour).UnixMilli(),
	}

	err = a.sessionRepo.Create(ctx, session)
	if err != nil {
		return "", err
	}

	return session.ID, nil
}

func isSamePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
