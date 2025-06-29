package config

import (
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type ObjectStorageConfig struct {
	BucketName string
	SecretKey  string
	AccessKey  string
	Endpoint   string
	Region     string
}

func LoadObjectStorageConfig() (*ObjectStorageConfig, error) {
	bucketName := util.GetEnvOr(common.OBJECT_STORAGE_BUCKET_KEY, "")
	secretKey := util.GetEnvOr(common.OBJECT_STORAGE_SECRET_KEY, "")
	accessKey := util.GetEnvOr(common.OBJECT_STORAGE_ACCESS_KEY, "")
	endpoint := util.GetEnvOr(common.OBJECT_STORAGE_ENDPOINT_KEY, "")
	region := util.GetEnvOr(common.OBJECT_STORAGE_REGION_KEY, "")

	if bucketName == "" ||
		secretKey == "" ||
		accessKey == "" ||
		endpoint == "" ||
		region == "" {
		return nil, fmt.Errorf("missing object storage config, secretKey=%s bucketName=%s accessKey=%s endpoint=%s region=%s: %w",
			secretKey, bucketName, accessKey, endpoint, region, common.ErrInternalServerError)
	}

	return &ObjectStorageConfig{
		BucketName: bucketName,
		SecretKey:  secretKey,
		AccessKey:  accessKey,
		Endpoint:   endpoint,
		Region:     region,
	}, nil
}
