package repo

import (
	"database/sql"

	"github.com/ahleongzc/leetcode-live-backend/internal/entity"
)

type SessionRepo interface {
	GetSessionByID(ID string) *entity.Session
	CreateSession(session *entity.Session) error
	DeleteSessionByID(ID string) *entity.Session
}

func NewSessionRepoImpl(
	db *sql.DB,
) SessionRepo {
	return &SessionRepoImpl{
		db: db,
	}
}

type SessionRepoImpl struct {
	db *sql.DB
}

// DeleteSessionByID implements SessionRepo.
func (s *SessionRepoImpl) DeleteSessionByID(ID string) *entity.Session {
	panic("unimplemented")
}

// CreateSession implements SessionRepo.
func (s *SessionRepoImpl) CreateSession(session *entity.Session) error {
	panic("unimplemented")
}

// GetSessionByID implements SessionRepo.
func (s *SessionRepoImpl) GetSessionByID(ID string) *entity.Session {
	panic("unimplemented")
}
