package service

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) CreateVariation(ctx context.Context, variation *domain.Variation) error {
	now := time.Now().UTC()
	variation.CreatedAt = &now
	variation.UpdatedAt = &now

	err := s.VariationRepo.Insert(ctx, variation)
	if err != nil {
		return sharedD.ManageError(err, "error creating variation")
	}

	return nil
}