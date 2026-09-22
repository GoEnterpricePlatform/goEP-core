package service

import (
	"context"
	"time"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
)

func (s *Service) Create(ctx context.Context, plan *domain.Plan) error {
	now := time.Now().UTC()
	plan.CreatedAt = &now
	plan.UpdatedAt = &now

	for _, pItem := range plan.Items {
		pItem.CreatedAt = &now
		pItem.UpdatedAt = &now
		pItem.Status = domain.PlanStatusActive
	}

	err := s.PlanRepo.Insert(ctx, plan)
	if err != nil {
		return sharedD.ManageError(err, "error inserting plan")
	}

	return nil
}
