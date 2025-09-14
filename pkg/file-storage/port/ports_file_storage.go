package port

import (
	"context"
	"io"
)

type FileStorageAdt interface {
	UploadImage(ctx context.Context, imgPath string, file io.Reader, contentType string) error
	GetImage(ctx context.Context, imgPath string) (string, error)
}

type FileStorageSrv interface {
	UploadImage(ctx context.Context, imgPath string, file io.Reader, contentType string) (string, error)
	GetImage(ctx context.Context, imgPath string) (string, error)
}
