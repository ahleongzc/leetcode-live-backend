package service

import (
	"context"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, sessionID string) error
	Register(ctx context.Context, email, password string) error
	ValidateSession(ctx context.Context, sessionID string) (bool, error)
}

func NewAuthService(
	sessionRepo repo.SessionRepo,
) AuthService {
	return &AuthServiceImpl{
		sessionRepo: sessionRepo,
	}
}

type AuthServiceImpl struct {
	sessionRepo repo.SessionRepo
}

// CheckIsLoggedIn implements AuthService.
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

// Logout implements AuthService.
func (a *AuthServiceImpl) Logout(ctx context.Context, sessionID string) error {
	err := a.sessionRepo.DeleteByID(ctx, sessionID)
	if err != nil {
		return err
	}

	return nil
}

// Register implements AuthService.
func (a *AuthServiceImpl) Register(ctx context.Context, email string, password string) error {
	panic("unimplemented")
}

// Login implements AuthService.
func (a *AuthServiceImpl) Login(ctx context.Context, email string, password string) (string, error) {
	panic("unimplemented")
}
