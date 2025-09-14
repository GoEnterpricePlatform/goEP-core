package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

func (a *Adapter) UploadImage(ctx context.Context, imgPath string, file io.Reader, contentType string) error {
	options := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := a.MinioClient.PutObject(ctx, a.BucketName, imgPath, file, -1, options)
	if err != nil {
		return fmt.Errorf("failed to upload image %q to bucket %q: %w", imgPath, a.BucketName, err)
	}

	return nil
}
