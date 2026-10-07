package service

import (
	"context"
	"math"

	paddlPlanD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	helpers "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/helper"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) GetAll(ctx context.Context, limit int64, page int64) ([]*paddlPlanD.PaddlePlan, int64, int64, error) {
	paddlePlans, err := s.PaddlePlanRepo.FindAll(ctx, limit, page)
	if err != nil {
		return nil, 0, 0, sharedD.ManageError(err, "error getting paddle plans")
	}

	for _, paddlePlan := range paddlePlans {
		plan, err := s.PlanRepo.Find(ctx, paddlePlan.PlanID)
		if err != nil {
			return nil, 0, 0, sharedD.ManageError(err, "error getting linked plan")
		}

		helpers.CalculateVariations(plan)

		paddlePlan.Name = plan.Name
		paddlePlan.Description = plan.Description
		paddlePlan.ImgPath = plan.ImgPath
		paddlePlan.Variations = plan.Variations

		if plan.ImgPath != nil {
			url, err := s.PaddlePlanStg.GetImage(ctx, *plan.ImgPath)
			if err != nil {
				return nil, 0, 0, sharedD.ManageError(err, "error getting plan image")
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
					return nil, 0, 0, sharedD.ManageError(
						err,
						"error getting plan item image",
					)
				}

				paddleItem.ImgUrl = &url
			}
		}

		s.enrichPaddleDetails(ctx, paddlePlan)
	}

	count, err := s.PaddlePlanRepo.Count(ctx)
	if err != nil {
		return nil, 0, 0, sharedD.ManageError(err, "error counting paddle plans")
	}

	totalPages := int64(math.Ceil(float64(count) / float64(limit)))

	return paddlePlans, count, totalPages, nil
}
