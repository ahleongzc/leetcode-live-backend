package entity

import (
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/util"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	Base
	Username             string
	Email                string `gorm:"index;unique"`
	Password             string
	LoginCount           uint
	LastLoginTimeStampMS *int64
	DeletedTimestampMS   soft_delete.DeletedAt `gorm:"softDelete:milli"`
	SettingID            uint
}

func NewUser() *User {
	return &User{}
}

func (u *User) SetSettingID(settingID uint) *User {
	if u == nil {
		return nil
	}
	u.SettingID = settingID
	return u
}

func (u *User) SetEmail(email string) *User {
	if u == nil {
		return nil
	}
	u.Email = email
	return u
}

func (u *User) SetPassword(password string) *User {
	if u == nil {
		return nil
	}
	u.Password = password
	return u
}

func (u *User) SetUsername(username string) *User {
	if u == nil {
		return nil
	}
	u.Username = username
	return u
}

func (u *User) SetLoginCount(count uint) *User {
	if u == nil {
		return u
	}
	u.LoginCount = count
	return u
}

func (u *User) SetLastLoginTimestampMS(timestampMS int64) *User {
	if u == nil {
		return nil
	}
	u.LastLoginTimeStampMS = util.ToPtr(timestampMS)

	return u
}

func (u *User) Login() {
	u.IncrementLoginCount()
	u.SetLastLoginTimestampMS(time.Now().UnixMilli())
}

func (u *User) IncrementLoginCount() {
	if u == nil {
		return
	}
	u.SetLoginCount(u.LoginCount + 1)
}
