package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID                uint   `gorm:"primarykey"`
	UUID              string `gorm:"type:string;size:36"`
	CreateTimestampMS int64  `gorm:"autoCreateTime:milli"`
	UpdateTimestampMS int64  `gorm:"autoUpdateTime:milli"`
}

// BeforeCreate will set a UUID if it's not already set
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	if base.UUID == "" {
		base.UUID = uuid.New().String()
	}
	return nil
}
