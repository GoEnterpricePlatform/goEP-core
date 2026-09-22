package core

import (
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)
 

func validatePlanFields(p CreateOrUpdatePlanReq) error {
	if strings.TrimSpace(p.Name) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "name is required")
	}

	if p.Description != nil && strings.TrimSpace(*p.Description) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "description is required")
	}

	// Plan must have at least one item.
	if len(p.PlanItems) == 0 {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "at least one plan item is required")
	}

	for _, item := range p.PlanItems {
		if item == nil {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "plan item cannot be null")
		}

		// Features validation.
		if len(item.Features) == 0 {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "at least one feature is required")
		}

		for _, feature := range item.Features {
			if strings.TrimSpace(feature) == "" {
				return domain.NewAppError(domain.ErrCodeInvalidParams, "feature cannot be empty")
			}
		}
	}

	// Plan with variants must have at least 2 items.
	if len(p.PlanItems) > 1 {
		for _, item := range p.PlanItems {
			if len(item.VarOptionIDs) == 0 {
				return domain.NewAppError(domain.ErrCodeInvalidParams, "variant options are required for plan items")
			}
		}
	}

	// Plan without variants can only have one item.
	if len(p.PlanItems) == 1 {
		item := p.PlanItems[0]

		if len(item.VarOptionIDs) > 0 {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "single plan must not contain variant options")
		}
	}

	return nil
}

func validatePlanItemIDs(p CreateOrUpdatePlanReq) error {
	for _, item := range p.PlanItems {
		if item == nil {
			continue
		}

		if strings.TrimSpace(item.ID) == "" {
			return domain.NewAppError(
				domain.ErrCodeInvalidParams,
				"plan item id is required",
			)
		}
	}

	return nil
}