package minio

import (
	"context"
	"fmt"
	"net/url"
)


func (fs *FileStorage) GetImage(ctx context.Context, imgPath string) (string, error) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)

	// Generates a presigned url which expires in a day.
	presignedURL, err := fs.MinioClient.PresignedGetObject(ctx, fs.BucketName, imgPath, fs.ExpTime, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL for user image %q in bucket %q: %w", imgPath, fs.BucketName, err)
	}

	return presignedURL.String(), nil
}

