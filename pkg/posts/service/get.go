package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Get(ctx context.Context, id string) (*domain.Post, error) {
	post, err := s.PostRepo.Find(ctx, id)
	if err != nil {
		return nil, sharedD.ManageError(err, "error getting post")
	}
	return post, nil
}
