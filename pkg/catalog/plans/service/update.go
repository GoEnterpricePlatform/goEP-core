package service

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Update(ctx context.Context, id string, plan *domain.Plan) error {
	now := time.Now().UTC()

	currentPlan, err := s.PlanRepo.Find(ctx, id)
	if err != nil {
		return sharedD.ManageError(err, "error finding plan")
	}

	// Preserve plan-specific data.
	plan.ID = currentPlan.ID
	plan.CreatedAt = currentPlan.CreatedAt
	plan.UpdatedAt = &now
	plan.ImgPath = currentPlan.ImgPath

	// Mix received items with existing ones.
	for _, item := range plan.Items {
		if item == nil {
			continue
		}

		var currentItem *domain.PlanItem

		for _, existingItem := range currentPlan.Items {
			if existingItem == nil {
				continue
			}

			if item.ID == existingItem.ID {
				currentItem = existingItem
				break
			}
		}

		// The item sent by the client does not belong to this plan.
		if currentItem == nil {
			return sharedD.NewAppError(
				sharedD.ErrCodeNotFound, "plan item with id "+item.ID+" does not exist",
			)
		}

		// Preserve fields that are not editable by this request.
		item.CreatedAt = currentItem.CreatedAt
		item.Status = currentItem.Status
		item.ImgPath = currentItem.ImgPath

		// Update timestamp because this item is being modified.
		item.UpdatedAt = &now
	}

	err = s.PlanRepo.Update(ctx, id, plan)
	if err != nil {
		return sharedD.ManageError(err, "error updating plan")
	}
	return nil
}
