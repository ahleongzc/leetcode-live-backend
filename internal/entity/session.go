package entity

type Session struct {
	Base
	Token             string `gorm:"index,unique"`
	UserID            uint
	ExpireTimestampMS int64
}
