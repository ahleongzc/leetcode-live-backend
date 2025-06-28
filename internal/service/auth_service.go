package service

import (
	"context"
	"errors"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, sessionID string) error
	ValidateSession(ctx context.Context, sessionID string) (bool, error)
}

func NewAuthService(
	sessionRepo repo.SessionRepo,
	userRepo repo.UserRepo,
) AuthService {
	return &AuthServiceImpl{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

type AuthServiceImpl struct {
	sessionRepo repo.SessionRepo
	userRepo    repo.UserRepo
}

func (a *AuthServiceImpl) ValidateSession(ctx context.Context, sessionID string) (bool, error) {
	session, err := a.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return false, err
	}

	if session.ExpireTimestampMS < time.Now().UnixMilli() {
		return false, nil
	}

	session.ExpireTimestampMS = time.Now().Add(48 * time.Hour).UnixMilli()
	err = a.sessionRepo.Update(ctx, session)
	if err != nil {
		return false, err
	}

	return true, nil
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

	return "", nil
}

func isSamePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
