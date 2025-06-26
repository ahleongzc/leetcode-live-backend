package service

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
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

// Login implements AuthService.
func (a *AuthServiceImpl) Login(ctx context.Context, email string, password string) (string, error) {
	panic("unimplemented")
}
