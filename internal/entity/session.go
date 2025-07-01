package entity

type Session struct {
	Base
	Token             string `gorm:"unique"`
	UserID            uint
	ExpireTimestampMS int64
}
