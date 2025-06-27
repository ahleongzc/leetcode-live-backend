package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
)

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
}

func NewUserRepo(database *sql.DB) UserRepo {
	return &UserRepoImpl{
		db: database,
	}
}

type UserRepoImpl struct {
	db *sql.DB
}

// GetByEmail implements UserRepo.
func (u *UserRepoImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{email}

	query := fmt.Sprintf(`
		SELECT 
			id, email, password, is_deleted, last_login_timestamp_ms 
		FROM 
			%s
		WHERE 
			email = $1
	`, common.USER_TABLE_NAME)

	user := &entity.User{}
	err := u.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Password, &user.IsDeleted, &user.LastLoginTimeStampMS)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get user with email %s, %s: %w", email, err, common.ErrInternalServerError)
	}

	return user, nil
}

// Create implements UserRepo.
func (u *UserRepoImpl) Create(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{user.Email, user.Password}

	query := fmt.Sprintf(`
		INSERT INTO %s
			(email, password)
		VALUES
			($1, $2)
	`, common.USER_TABLE_NAME)

	_, err := u.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to create new user, %s: %w", err, common.ErrInternalServerError)
	}

	return nil
}

// Delete implements UserRepo.
func (u *UserRepoImpl) Delete(ctx context.Context, id int) error {
	panic("unimplemented")
}

// GetByID implements UserRepo.
func (u *UserRepoImpl) GetByID(ctx context.Context, id int) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{id}

	query := fmt.Sprintf(`
		SELECT 
			id, email, password, is_deleted, last_login_timestamp_ms 
		FROM 
			%s
		WHERE 
			id = $1
	`, common.USER_TABLE_NAME)

	user := &entity.User{}
	err := u.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Password, &user.IsDeleted, &user.LastLoginTimeStampMS)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user: %w", common.ErrNotFound)
		}
		return nil, fmt.Errorf("unable to get user with id %d, %s: %w", id, err.Error(), common.ErrInternalServerError)
	}

	return user, nil
}

// Update implements UserRepo.
func (u *UserRepoImpl) Update(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, common.DB_QUERY_TIMEOUT)
	defer cancel()

	args := []any{user.Email, user.Password, user.IsDeleted, user.LastLoginTimeStampMS, user.ID}

	query := fmt.Sprintf(`
		UPDATE 
			%s
		SET 
			email = $1,
		    password = $2,
		    is_deleted = $3,
		    last_login_timestamp_ms = $4
		WHERE 
			id = $5
	`, common.USER_TABLE_NAME)

	result, err := u.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to update user, %s: %w", err, common.ErrInternalServerError)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to update user, %s: %w", err, common.ErrInternalServerError)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("unable to update user, affected row count is 0: %w", common.ErrInternalServerError)
	}

	return nil
}
