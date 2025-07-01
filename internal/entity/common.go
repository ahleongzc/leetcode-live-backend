package entity

type Base struct {
	ID                uint  `gorm:"primarykey"`
	CreateTimestampMS int64 `gorm:"autoCreateTime:milli"`
	UpdateTimestampMS int64 `gorm:"autoUpdateTime:milli"`
}
