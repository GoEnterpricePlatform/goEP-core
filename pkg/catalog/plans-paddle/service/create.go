package service

import (
	"context"
	"time"

	paddlePlanD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	planD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Create(ctx context.Context, paddlePlan *paddlePlanD.PaddlePlan) error {
	now := time.Now().UTC()

	paddlePlan.CreatedAt = &now
	paddlePlan.UpdatedAt = &now

	plan := &planD.Plan{
		Name:        paddlePlan.Name,
		Description: paddlePlan.Description,
		Items:       make([]*planD.PlanItem, 0, len(paddlePlan.Items)),
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	for _, paddleItem := range paddlePlan.Items {
		if paddleItem == nil {
			continue
		}
		plan.Items = append(plan.Items, &planD.PlanItem{
			Features:     paddleItem.Features,
			VarOptionIDs: paddleItem.VarOptionIDs,
			Status:       planD.PlanStatusActive,
			CreatedAt:    &now,
			UpdatedAt:    &now,
		})

		paddleItem.CreatedAt = &now
		paddleItem.UpdatedAt = &now
	}

	if err := s.PaddlePlanTx.Create(ctx, plan, paddlePlan); err != nil {
		return sharedD.ManageError(err, "error creating paddle plan")
	}

	for i, planItem := range plan.Items {
		if planItem == nil || i >= len(paddlePlan.Items) || paddlePlan.Items[i] == nil {
			continue
		}

		paddlePlan.Items[i].PlanItemID = planItem.ID
		paddlePlan.Items[i].Features = planItem.Features
		paddlePlan.Items[i].VarOptionIDs = planItem.VarOptionIDs
		paddlePlan.Items[i].Status = planItem.Status
	}

	return nil
}
