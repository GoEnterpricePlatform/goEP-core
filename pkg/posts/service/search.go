package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Search(ctx context.Context, query string, limit int) ([]*domain.Post, error) {
	posts, err := s.PostRepo.Search(ctx, query, limit)
	if err != nil {
		return nil, sharedD.ManageError(err, "error searching posts")
	}
	return posts, nil
}
