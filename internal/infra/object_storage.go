package infra

import (
	"context"
	"fmt"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	internalConfig "github.com/ahleongzc/leetcode-live-backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewCloudflareR2ObjectStorageClient(
	objectStorageConfig *internalConfig.ObjectStorageConfig,
) (*s3.Client, error) {
	if objectStorageConfig == nil {
		return nil, fmt.Errorf("no object storage config when initialising s3 client: %w", common.ErrInternalServerError)
	}

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				objectStorageConfig.AccessKey,
				objectStorageConfig.SecretKey,
				"",
			),
		),
		config.WithRegion("auto"),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(objectStorageConfig.Endpoint)
	})

	if client == nil {
		return nil, fmt.Errorf("client is nil after initialization from config: %w", common.ErrInternalServerError)
	}

	return client, err
}
