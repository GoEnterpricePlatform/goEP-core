package service

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) UpdateVariation(ctx context.Context, variation *domain.Variation) error {
	now := time.Now().UTC()
	variation.UpdatedAt = &now

	err := s.VariationRepo.Update(ctx, variation)
	if err != nil {
		return sharedD.ManageError(err, "error updating variation")
	}
	return nil
}
