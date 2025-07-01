package entity

import "gorm.io/plugin/soft_delete"

type User struct {
	Base
	Email                string
	Password             string
	LoginCount           int
	LastLoginTimeStampMS *int64
	DeletedTimestampMS   soft_delete.DeletedAt `gorm:"softDelete:milli"`
}
