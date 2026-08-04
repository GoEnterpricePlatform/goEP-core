package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) GetAllVariationsWithOptions(ctx context.Context) ([]*domain.Variation, error) {
	variations, err := s.VariationRepo.FindAllWithOptions(ctx)
	if err != nil {
		return nil, sharedD.ManageError(err, "error getting variations with options")
	}
    return variations, nil
}