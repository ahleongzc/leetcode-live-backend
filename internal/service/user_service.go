package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterNewUser(ctx context.Context, email, password string) error
}

func NewUserService(
	userRepo repo.UserRepo,
) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

type UserServiceImpl struct {
	userRepo repo.UserRepo
}

// RegisterNewUser implements UserService.
func (u *UserServiceImpl) RegisterNewUser(ctx context.Context, email, password string) error {
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format: %w", common.ErrBadRequest)
	}

	if len(password) > 20 {
		return fmt.Errorf("password is too long, must be less than or equal to 20 characters: %w", common.ErrBadRequest)
	}

	if len(password) < 8 {
		return fmt.Errorf("password is too short, must be longer than or equal to 8 characters: %w", common.ErrBadRequest)
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

	user := &entity.User{
		Email:    email,
		Password: hashedPassword,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
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
