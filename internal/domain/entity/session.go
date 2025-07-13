package entity

import "time"

type Session struct {
	Base
	Token             string `gorm:"index"`
	UserID            uint
	ExpireTimestampMS int64
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) SetToken(token string) *Session {
	if s == nil {
		return nil
	}
	s.Token = token
	return s
}

func (s *Session) SetUserID(userID uint) *Session {
	if s == nil {
		return nil
	}
	s.UserID = userID
	return s
}

func (s *Session) IsExpired() bool {
	if s == nil {
		return true
	}

	return s.ExpireTimestampMS < time.Now().UnixMilli()
}

func (s *Session) SetExpireTimestampMS(expireTimestampMS int64) *Session {
	if s == nil {
		return nil
	}
	s.ExpireTimestampMS = expireTimestampMS
	return s
}

func (s *Session) SetExpireTimestampUsingDays(dayCount uint) *Session {
	if s == nil {
		return nil
	}
	hours := time.Duration(dayCount) * 24 * time.Hour
	s.SetExpireTimestampMS(time.Now().Add(hours).UnixMilli())
	return s
}

func (s *Session) AddDayCountToPreviousExpireTimestampMS(dayCount uint) {
	if s == nil {
		return
	}
	hours := time.Duration(dayCount) * 24 * time.Hour
	s.SetExpireTimestampMS(time.UnixMilli(s.ExpireTimestampMS).Add(hours).UnixMilli())
}
