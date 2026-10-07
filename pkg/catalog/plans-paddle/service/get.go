package service

import (
	"context"

	paddlPlanD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	helpers "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/helper"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// tambien las imagenes si no es nil ImgPath
func (s *Service) Get(ctx context.Context, id string) (*paddlPlanD.PaddlePlan, error) {
	paddlePlan, err := s.PaddlePlanRepo.Find(ctx, id)
	if err != nil {
		return nil, sharedD.ManageError(err, "error getting paddle plan")
	}

	plan, err := s.PlanRepo.Find(ctx, paddlePlan.PlanID)
	if err != nil {
		return nil, sharedD.ManageError(err, "error getting plan")
	}

	helpers.CalculateVariations(plan)

	paddlePlan.Name = plan.Name
	paddlePlan.Description = plan.Description
	paddlePlan.ImgPath = plan.ImgPath
	paddlePlan.Variations = plan.Variations

	if plan.ImgPath != nil {
		url, err := s.PaddlePlanStg.GetImage(ctx, *plan.ImgPath)
		if err != nil {
			return nil, sharedD.ManageError(err, "error getting plan image")
		}

		paddlePlan.ImgUrl = &url
	}

	planItems := make(map[string]*domain.PlanItem, len(plan.Items))

	for _, item := range plan.Items {
		planItems[item.ID] = item
	}

	for _, paddleItem := range paddlePlan.Items {
		planItem, ok := planItems[paddleItem.PlanItemID]
		if !ok {
			continue
		}

		paddleItem.ImgPath = planItem.ImgPath
		paddleItem.Features = planItem.Features
		paddleItem.Status = planItem.Status
		paddleItem.VarOptionIDs = planItem.VarOptionIDs
		paddleItem.Options = planItem.Options

		if planItem.ImgPath != nil {
			url, err := s.PaddlePlanStg.GetImage(ctx, *planItem.ImgPath)
			if err != nil {
				return nil, sharedD.ManageError(err, "error getting plan item image")
			}

			paddleItem.ImgUrl = &url
		}
	}

	return paddlePlan, nil
}
