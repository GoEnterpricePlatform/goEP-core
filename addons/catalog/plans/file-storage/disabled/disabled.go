package disabled

import (
	"context"
	"errors"
	"io"
)

type FileStorage struct{}

func NewPlanDisabledAdapter() *FileStorage {
	return &FileStorage{}
}

func (fs *FileStorage) GetImage(ctx context.Context, imgPath string) (string, error) {
	return "", errors.New("error plan file storage disabled")
}


func (fs *FileStorage) UploadImage(ctx context.Context, imgPath string, file io.Reader, contentType string) error {
	return errors.New("error plan file storage disabled")
}
