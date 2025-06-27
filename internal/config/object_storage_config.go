package config

import (
	"os"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

type ObjectStorageConfig struct {
	BucketName string
	SecretKey  string
	AccessKey  string
	Endpoint   string
	Region     string
}

func LoadObjectStorageConfig() *ObjectStorageConfig {
	bucketName := os.Getenv(common.R2_BUCKET_KEY)
	secretKey := os.Getenv(common.R2_SECRET_KEY)
	accessKey := os.Getenv(common.R2_ACCESS_KEY)
	endpoint := os.Getenv(common.R2_ENDPOINT_KEY)
	region := os.Getenv(common.R2_REGION_KEY)

	return &ObjectStorageConfig{
		BucketName: bucketName,
		SecretKey:  secretKey,
		AccessKey:  accessKey,
		Endpoint:   endpoint,
		Region:     region,
	}
}
