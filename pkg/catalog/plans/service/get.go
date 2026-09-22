package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Get(ctx context.Context, id string) (*domain.Plan, error) {
	plan, err := s.PlanRepo.Find(ctx, id)
	if err != nil {
		return nil, sharedD.ManageError(err, "error getting plan")
	}
	return plan, nil
}
