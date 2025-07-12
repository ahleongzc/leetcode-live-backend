package entity

import "gorm.io/plugin/soft_delete"

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

func NewUserWithSettingID(settingID uint) *User {
	return &User{
		SettingID: settingID,
	}
}

func (u *User) SetEmail(email string) {
	if u == nil {
		return
	}
	u.Email = email
}

func (u *User) SetPassword(password string) {
	if u == nil {
		return
	}
	u.Password = password
}

func (u *User) SetUsername(username string) {
	if u == nil {
		return
	}
	u.Username = username
}
