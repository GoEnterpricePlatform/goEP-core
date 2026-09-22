package service

import (
	"context"
	"math"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	helpers "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/helper"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) GetAll(ctx context.Context, limit int64, page int64) ([]*domain.Plan, int64, int64, error) {
	plans, err := s.PlanRepo.FindAll(ctx, limit, page)
	if err != nil {
		return nil, 0, 0, sharedD.ManageError(err, "error getting plans")
	}

	for _, plan := range plans {
		if plan.ImgPath != nil {
			url, err := s.PlanFileStg.GetImage(ctx, *plan.ImgPath)
			if err != nil {
				return nil, 0, 0, sharedD.ManageError(err, "")
			}
			plan.ImgUrl = &url
		}

		for _, pItem := range plan.Items {
			// esto depeden si es para el admin o es para el user
			// revisar para esta validacion para otros handler como get
			/* if pItem.Status != domain.PlanStatusActive {
				continue
			} */

			if pItem.ImgPath != nil {
				url, err := s.PlanFileStg.GetImage(ctx, *pItem.ImgPath)
				if err != nil {
					return nil, 0, 0, sharedD.ManageError(err, "")
				}
				pItem.ImgUrl = &url
			}

			pItem.VarOptionIDs = nil
		}

		helpers.CalculateVariations(plan)
	}

	count, err := s.PlanRepo.Count(ctx)
	if err != nil {
		return nil, 0, 0, sharedD.ManageError(err, "")
	}
	totalPages := int64(math.Ceil(float64(count) / float64(limit)))

	return plans, count, totalPages, nil
}
