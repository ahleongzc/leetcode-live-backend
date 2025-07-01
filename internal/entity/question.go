package entity

type Question struct {
	Base
	ExternalID  string `gorm:"index"`
	Description string
}
