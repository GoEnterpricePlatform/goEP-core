package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) UploadAvatar(ctx context.Context, img string, file io.Reader, userID string, contentType string) error {
	exists, err := s.UserRepo.Exists(ctx, userID)
	if err != nil {
		return domain.ManageError(err, "")
	}

	if !exists {
		return domain.ManageError(domain.ErrNotFound, "")
	}

	path := fmt.Sprintf("%s_users/%s%s", s.MdlName, userID, filepath.Ext(img))

	err = s.UserFileStg.UploadImage(ctx, path, file, contentType)
	if err != nil {
		return domain.ManageError(err, "")
	}

	now := time.Now().UTC()
	err = s.UserRepo.UpdateAvatarPath(ctx, userID, path, now)
	if err != nil {
		return domain.ManageError(err, "")
	}

	return nil
}
