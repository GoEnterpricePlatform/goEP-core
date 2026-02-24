package service

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Create(ctx context.Context, post *domain.Post) error {
	now := time.Now().UTC()
	post.CreatedAt = &now

	err := s.PostRepo.Insert(ctx, post)
	if err != nil {
		return sharedD.ManageError(err, "")
	}

	return nil
}
