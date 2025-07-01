package entity

type Base struct {
	ID                 uint  `gorm:"primarykey"`
	CreatedTimestampMS int64 `gorm:"autoCreateTime:milli"`
	UpdatedTimestampMS int64 `gorm:"autoUpdateTime:milli"`
}
