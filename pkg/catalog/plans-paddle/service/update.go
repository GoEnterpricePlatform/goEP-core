package service

import (
	"context"
	"time"

	paddleD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	planD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Update(ctx context.Context, id string, paddlePlan *paddleD.PaddlePlan) error {
	current, err := s.PaddlePlanRepo.Find(ctx, id)
	if err != nil {
		return sharedD.ManageError(err, "error finding paddle plan")
	}
	currentPlan, err := s.PlanRepo.Find(ctx, current.PlanID)
	if err != nil {
		return sharedD.ManageError(err, "error finding plan")
	}
	planItems := make(map[string]*planD.PlanItem, len(currentPlan.Items))
	for _, item := range currentPlan.Items {
		if item == nil {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "linked plan contains a missing item")
		}
		planItems[item.ID] = item
	}
	paddleItems := make(map[string]*paddleD.PaddlePlanItem, len(current.Items))
	for _, item := range current.Items {
		if item == nil {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "paddle plan contains a missing item")
		}
		paddleItems[item.ID] = item
	}
	now := time.Now().UTC()
	plan := &planD.Plan{ID: currentPlan.ID, Name: paddlePlan.Name, Description: paddlePlan.Description, ImgPath: currentPlan.ImgPath, CreatedAt: currentPlan.CreatedAt, UpdatedAt: &now, Items: make([]*planD.PlanItem, 0, len(paddlePlan.Items))}
	seen := make(map[string]bool, len(paddlePlan.Items))
	for _, item := range paddlePlan.Items {
		if item == nil || (item.ID != "" && seen[item.ID]) {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "duplicate or missing plan item")
		}
		if item.ID == "" {
			item.CreatedAt = &now
			item.UpdatedAt = &now
			plan.Items = append(plan.Items, &planD.PlanItem{Features: item.Features, VarOptionIDs: item.VarOptionIDs, Status: planD.PlanStatusActive, CreatedAt: &now, UpdatedAt: &now})
			continue
		}
		seen[item.ID] = true
		existing, ok := paddleItems[item.ID]
		if !ok {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "plan item does not belong to paddle plan")
		}
		base, ok := planItems[existing.PlanItemID]
		if !ok {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "linked plan item does not exist")
		}
		item.PlanItemID = existing.PlanItemID
		item.CreatedAt = existing.CreatedAt
		item.UpdatedAt = &now
		plan.Items = append(plan.Items, &planD.PlanItem{ID: base.ID, Features: item.Features, VarOptionIDs: item.VarOptionIDs, Status: base.Status, ImgPath: base.ImgPath, CreatedAt: base.CreatedAt, UpdatedAt: &now})
	}
	paddlePlan.ID = current.ID
	paddlePlan.PlanID = current.PlanID
	paddlePlan.CreatedAt = current.CreatedAt
	paddlePlan.UpdatedAt = &now
	if err := s.PaddlePlanTx.Update(ctx, plan, paddlePlan); err != nil {
		return sharedD.ManageError(err, "error updating paddle plan")
	}
	return nil
}
