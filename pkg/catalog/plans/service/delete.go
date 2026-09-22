package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.PlanRepo.Delete(ctx, id)
	if err != nil {
		return domain.ManageError(err, "error deleting plan")
	}
	return nil
}