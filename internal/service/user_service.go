package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterNewUser(ctx context.Context, email, password string) error
	GetUserProfile(ctx context.Context, userID uint) (*model.UserProfile, error)
}

func NewUserService(
	userRepo repo.UserRepo,
	settingRepo repo.SettingRepo,
) UserService {
	return &UserServiceImpl{
		userRepo:    userRepo,
		settingRepo: settingRepo,
	}
}

type UserServiceImpl struct {
	userRepo    repo.UserRepo
	settingRepo repo.SettingRepo
}

// GetUserProfile implements UserService.
func (u *UserServiceImpl) GetUserProfile(ctx context.Context, userID uint) (*model.UserProfile, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	setting, err := u.settingRepo.GetByID(ctx, user.SettingID)
	if err != nil {
		return nil, err
	}

	userProfile := model.NewUserProfile().
		SetEmail(user.Email).
		SetUsername(user.Username).
		SetRemainingInterviewCount(setting.RemainingInterviewCount).
		SetInterviewDurationS(uint(setting.InterviewDurationS))

	return userProfile, nil
}

// RegisterNewUser implements UserService.
func (u *UserServiceImpl) RegisterNewUser(ctx context.Context, email, password string) error {
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format: %w", common.ErrBadRequest)
	}

	if !isValidPassword(password) {
		return fmt.Errorf("invalid password length, must be between 8 and 20 characters: %w", common.ErrBadRequest)
	}

	isDuplicated, err := u.isDuplicatedEmail(ctx, email)
	if err != nil {
		return err
	}

	if isDuplicated {
		return fmt.Errorf("email is already taken: %w", common.ErrBadRequest)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("unable to hash password, %s: %w", err, common.ErrInternalServerError)
	}

	settingID, err := u.createDefaultSetting(ctx)
	if err != nil {
		return nil
	}

	user := entity.NewUser().
		SetSettingID(settingID).
		SetEmail(email).
		SetPassword(hashedPassword).
		SetUsername(email)

	if err := u.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) createDefaultSetting(ctx context.Context) (uint, error) {
	setting := entity.NewDefaultSetting()

	settingID, err := u.settingRepo.Create(ctx, setting)
	if err != nil {
		return 0, err
	}

	return settingID, nil
}

func isValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 20
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UserServiceImpl) isDuplicatedEmail(ctx context.Context, email string) (bool, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return false, nil
		}
		return true, err
	}

	return user != nil, nil
}
