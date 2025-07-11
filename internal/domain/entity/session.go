package entity

type Session struct {
	Base
	Token             string `gorm:"index"`
	UserID            uint
	ExpireTimestampMS int64
}
