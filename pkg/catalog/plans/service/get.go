package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	helpers "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/helper"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Get(ctx context.Context, id string) (*domain.Plan, error) {
	plan, err := s.PlanRepo.Find(ctx, id)
	if err != nil {
		return nil, sharedD.ManageError(err, "error getting plan")
	}

	if plan.ImgPath != nil {
		url, err := s.PlanFileStg.GetImage(ctx, *plan.ImgPath)
		if err != nil {
			return nil, sharedD.ManageError(err, "error getting plan image")
		}

		plan.ImgUrl = &url
	}

	for _, item := range plan.Items {
		if item.ImgPath != nil {
			url, err := s.PlanFileStg.GetImage(ctx, *item.ImgPath)
			if err != nil {
				return nil, sharedD.ManageError(err, "error getting plan item image")
			}

			item.ImgUrl = &url
		}

		item.VarOptionIDs = nil
	}

	helpers.CalculateVariations(plan)

	return plan, nil
}
