package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) DeleteVarOption(ctx context.Context, id string, variationID string) error {
	inUse, err := s.OptionUsage.IsOptionInUse(ctx, id)
	if err != nil {
		return domain.ManageError(err, "error checking variation option usage")
	}
	if inUse {
		return domain.NewAppError(domain.ErrCodeConflict, "Cannot delete variation option while it is in use.")
	}
	err = s.VarOptionRepo.Delete(ctx, id, variationID)
	if err != nil {
		return domain.ManageError(err, "error deleting varOption")
	}
	return nil
}
