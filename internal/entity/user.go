package entity

type User struct {
	ID                   int
	Email                string
	Password             string
	LoginCount           int
	IsDeleted            bool
	LastLoginTimeStampMS *int64
}
