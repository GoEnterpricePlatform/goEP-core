package service

import (
	"context"

	paddlePlanD "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	planD "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)

func (s *Service) ResolvePaddleCheckout(ctx context.Context, paddlePlanID, paddlePlanItemID string) (*paddlePlanD.PaddleCheckoutOffer, error) {
	if paddlePlanID == "" || paddlePlanItemID == "" {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "Paddle plan and plan item are required")
	}
	paddlePlan, err := s.PaddlePlanRepo.Find(ctx, paddlePlanID)
	if err != nil {
		return nil, sharedD.ManageError(err, "error finding Paddle plan for checkout")
	}

	var selected *paddlePlanD.PaddlePlanItem
	for _, item := range paddlePlan.Items {
		if item != nil && item.ID == paddlePlanItemID {
			selected = item
			break
		}
	}
	if selected == nil {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "Paddle plan item not found")
	}
	if selected.PaddlePriceID == "" {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "Paddle plan item has no Paddle price")
	}

	catalogPlan, err := s.PlanRepo.Find(ctx, paddlePlan.PlanID)
	if err != nil {
		return nil, sharedD.ManageError(err, "error finding Catalog plan for checkout")
	}
	for _, item := range catalogPlan.Items {
		if item == nil || item.ID != selected.PlanItemID {
			continue
		}
		if item.Status != planD.PlanStatusActive {
			return nil, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "plan item is not available for purchase")
		}
		return &paddlePlanD.PaddleCheckoutOffer{
			PaddlePriceID: selected.PaddlePriceID,
			PlanID:        paddlePlan.PlanID,
			PlanItemID:    selected.PlanItemID,
		}, nil
	}
	return nil, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "plan item is not available for purchase")
}
