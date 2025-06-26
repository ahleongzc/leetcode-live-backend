package entity

type Session struct {
	ID                string
	UserID            int
	ExpireTimestampMS int64
}
