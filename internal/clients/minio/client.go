package minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client
}

func NewClient(endPoint string, accessKeyID string, secretKeyID string, useSSL bool) (*MinioClient, error) {

	minioClient, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretKeyID, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client %w", err)
	}

	client := &MinioClient{
		Client: minioClient,
	}

	return client, nil
}

func (client *MinioClient) CreateStorage(bucketName string) error {
	found, err := client.Client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket %q exists: %w", bucketName, err)
	}

	if !found {
		err := client.Client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket %q: %w", bucketName, err)
		}
	}
	return nil
}
