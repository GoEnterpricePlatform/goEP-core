package service

import (
	"context"

	"github.com/amorindev/go-cms-tmpl/pkg/posts/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (s *Service) GetAll(ctx context.Context) ([]*domain.Post, error) {
	posts, err := s.PostRepo.FindAll(ctx)
	if err != nil {
		return nil, sharedD.ManageError(err, "")
	}
	return posts, nil
}