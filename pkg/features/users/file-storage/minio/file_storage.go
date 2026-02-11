package minio

import (
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/features/users/port"
	"github.com/minio/minio-go/v7"
)

var _ port.UserFileStg = &FileStorage{}

type FileStorage struct {
	MinioClient *minio.Client
	BucketName  string
	ExpTime     time.Duration
}

func NewUserFileStg(client *minio.Client, bucketName string, expTime time.Duration) *FileStorage {
	if expTime == 0 {
		expTime = time.Hour * 24 * 7
	}
	return &FileStorage{
		MinioClient: client,
		BucketName:  bucketName,
		ExpTime:     expTime,
	}
}
