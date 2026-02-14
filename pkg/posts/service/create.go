package service

import (
	"context"
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/posts/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
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
