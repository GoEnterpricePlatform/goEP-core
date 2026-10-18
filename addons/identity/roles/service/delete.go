package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.RoleRepo.Delete(ctx, id)
	if err != nil {
		return domain.ManageError(err, "error eliminating role")
	}
	return nil
}

