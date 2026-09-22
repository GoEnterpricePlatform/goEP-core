package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Patch(ctx context.Context, id string, plan *domain.Plan) (*domain.Plan, error) {
	existing, err := s.PlanRepo.Find(ctx, id)
	if err != nil {
		return nil, sharedD.ManageError(err, "")
	}

	// We only update the submitted fields.
	if plan.Name != "" {
		existing.Name = plan.Name
	}

	if plan.Description != nil {
		existing.Description = plan.Description
	}

	if plan.Items != nil {
		for _, item := range plan.Items {

			itemFound := false

			for _, existingItem := range existing.Items {

				if existingItem == nil {
					continue
				}

				if existingItem.ID != item.ID {
					continue
				}

				itemFound = true

				if item.Features != nil {
					existingItem.Features = item.Features
				}

				if item.VarOptionIDs != nil {
					existingItem.VarOptionIDs = item.VarOptionIDs
				}

				break
			}

			if !itemFound {
				errMsg := fmt.Sprintf("plan item with id %s not found", item.ID)
				return nil, sharedD.ManageError(sharedD.ErrNotFound, errMsg)

			}

		}
	}

	now := time.Now().UTC()
	existing.UpdatedAt = &now

	if err := s.PlanRepo.Update(ctx, id, existing); err != nil {
		return nil, sharedD.ManageError(err, "")
	}

	return existing, nil
}
