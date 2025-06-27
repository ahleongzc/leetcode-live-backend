package repo

import (
	"context"
	"fmt"
	"io"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileRepo interface {
	Upload(ctx context.Context, name string, content io.Reader, metadata map[string]any) (string, error)
}

// TODO: Move the client to an interface
func NewFileRepo(
	client *s3.Client,
	objectStorageConfig *config.ObjectStorageConfig,
) FileRepo {
	return &FileRepoImpl{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    objectStorageConfig.BucketName,
	}
}

type FileRepoImpl struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
}

func (f *FileRepoImpl) Upload(ctx context.Context, name string, content io.Reader, metadata map[string]any) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, common.R2_FILE_UPLOAD_TIMEOUT)
	defer cancel()

	res, err := f.presignClient.PresignPutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(f.bucketName),
			Key:    aws.String("example.txt"),
		},
		s3.WithPresignExpires(common.R2_FILE_UPLOAD_TIMEOUT),
	)

	if err != nil {
		return "", fmt.Errorf("unable to upload file to object storage, %s: %w", err, common.ErrInternalServerError)
	}

	return res.URL, nil
}
