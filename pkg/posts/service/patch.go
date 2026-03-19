package service

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Patch(ctx context.Context, id string, post *domain.Post) (*domain.Post, error) {
	existing, err := s.PostRepo.Find(ctx, id)
	if err != nil {
		return nil, sharedD.ManageError(err, "")
	}

	// We only update the submitted fields
	if post.Title != "" {
		existing.Title = post.Title
	}

	if post.Desc != nil {
		existing.Desc = post.Desc
	}

	now := time.Now().UTC()
	existing.UpdatedAt = &now

	err = s.PostRepo.Update(ctx, id, existing)
	if err != nil {
		return nil, sharedD.ManageError(err, "")
	}

	return existing, nil
}
