package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"
)

func (s *Service) UploadImage(ctx context.Context, imgPath string, file io.Reader, contentType string) (string, error) {
	fileName := imgPath[:len(imgPath)-len(filepath.Ext(imgPath))]
	uniqueFileName := fmt.Sprintf("%s-%d%s", fileName, time.Now().UnixNano(), filepath.Ext(imgPath))

	err := s.FileStgAdt.UploadImage(ctx, uniqueFileName, file, contentType)
	if err != nil {
		return "", err
	}
	return uniqueFileName, nil
}
