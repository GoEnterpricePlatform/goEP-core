package service

import (
	"context"

	"github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.PostRepo.Delete(ctx, id)
	if err != nil {
		return domain.ManageError(err, "")
	}
	return nil
}
