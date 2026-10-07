package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) DeleteVariation(ctx context.Context, id string) error {
	hasOptions, err := s.VarOptionRepo.HasByVariation(ctx, id)
	if err != nil {
		return domain.ManageError(err, "error checking variation options")
	}
	if hasOptions {
		return domain.NewAppError(domain.ErrCodeConflict, "Cannot delete variation while it has options. Delete its options first.")
	}
	err = s.VariationRepo.Delete(ctx, id)
	if err != nil {
		return domain.ManageError(err, "error deleting variation")
	}
	return nil
}
