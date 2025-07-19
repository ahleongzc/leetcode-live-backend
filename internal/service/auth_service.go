package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, token string) error
	GenerateRandomToken() string
	// Returns the interview ID
	ValidateAndConsumeInterviewToken(ctx context.Context, interviewToken string) (uint, error)
	// These two functions are only called by middleware
	ValidateAndRefreshSessionToken(ctx context.Context, token string) (string, error)
	GetUserIDFromSessionToken(ctx context.Context, token string) (uint, error)
}

func NewAuthService(
	userRepo repo.UserRepo,
	sessionRepo repo.SessionRepo,
	interviewRepo repo.InterviewRepo,
) AuthService {
	return &AuthServiceImpl{
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		interviewRepo: interviewRepo,
	}
}

type AuthServiceImpl struct {
	userRepo      repo.UserRepo
	sessionRepo   repo.SessionRepo
	interviewRepo repo.InterviewRepo
}

// ValidateAndConsumeInterviewToken implements AuthService.
func (a *AuthServiceImpl) ValidateAndConsumeInterviewToken(ctx context.Context, interviewToken string) (uint, error) {
	interview, err := a.interviewRepo.GetByToken(ctx, interviewToken)
	if err != nil {
		return 0, err
	}

	if !interview.Exists() {
		return 0, common.ErrUnauthorized
	}

	interview.ConsumeToken()

	err = a.interviewRepo.Update(ctx, interview)
	if err != nil {
		return 0, err
	}

	return interview.ID, nil
}

func (a *AuthServiceImpl) GenerateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return uuid.NewString()
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// GetUserIDFromSessionToken implements AuthService.
func (a *AuthServiceImpl) GetUserIDFromSessionToken(ctx context.Context, token string) (uint, error) {
	session, err := a.sessionRepo.GetByToken(ctx, token)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return 0, err
	}

	if !session.Exists() {
		return 0, common.ErrUnauthorized
	}

	user, err := a.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// ValidateSession implements AuthService.
func (a *AuthServiceImpl) ValidateAndRefreshSessionToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", common.ErrUnauthorized
	}

	session, err := a.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return "", common.ErrUnauthorized
		}
		return "", err
	}

	if session.IsExpired() {
		return "", common.ErrUnauthorized
	}

	session.AddDayCountToPreviousExpireTimestampMS(2)

	if err := a.sessionRepo.Update(ctx, session); err != nil {
		return "", err
	}

	return session.Token, nil
}

func (a *AuthServiceImpl) Logout(ctx context.Context, token string) error {
	err := a.sessionRepo.DeleteByToken(ctx, token)
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

	user.Login()

	if err := a.userRepo.Update(ctx, user); err != nil {
		return "", err
	}

	session := entity.NewSession().
		SetToken(a.GenerateRandomToken()).
		SetUserID(user.ID).
		SetExpireTimestampUsingDays(2)

	err = a.sessionRepo.Create(ctx, session)
	if err != nil {
		return "", err
	}

	return session.Token, nil
}

func isSamePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
