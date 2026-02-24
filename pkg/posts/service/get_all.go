package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) GetAll(ctx context.Context) ([]*domain.Post, error) {
	posts, err := s.PostRepo.FindAll(ctx)
	if err != nil {
		return nil, sharedD.ManageError(err, "")
	}
	return posts, nil
}
