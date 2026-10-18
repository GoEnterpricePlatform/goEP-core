package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

func (fs *FileStorage) UploadImage(ctx context.Context, imgPath string, file io.Reader, contentType string) error {
	options := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := fs.MinioClient.PutObject(ctx, fs.BucketName, imgPath, file, -1, options)
	if err != nil {
		return fmt.Errorf("failed to upload user image %q to bucket %q: %w", imgPath, fs.BucketName, err)
	}

	return nil
}
