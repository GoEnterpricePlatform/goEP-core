package core

import (
	"fmt"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (r CreatePaddlePlanReq) Validate() error {
	if r.Order < 0 {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "order cannot be negative")
	}

	if strings.TrimSpace(r.Name) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "name is required")
	}

	if r.Description != nil && strings.TrimSpace(*r.Description) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "description cannot be blank")
	}

	if strings.TrimSpace(r.PaddleProductID) == "" {
		return domain.NewAppError(
			domain.ErrCodeInvalidParams,
			"paddle product id is required for plan",
		)
	}

	if len(r.Items) == 0 {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "at least one plan item is required")
	}

	for i, item := range r.Items {
		if item == nil {
			return domain.NewAppError(
				domain.ErrCodeInvalidParams,
				"plan item at index "+fmt.Sprint(i)+" is required",
			)
		}

		if len(item.Features) == 0 {
			return domain.NewAppError(
				domain.ErrCodeInvalidParams,
				"plan item at index "+fmt.Sprint(i)+" requires at least one feature",
			)
		}

		for _, feature := range item.Features {
			if strings.TrimSpace(feature) == "" {
				return domain.NewAppError(
					domain.ErrCodeInvalidParams,
					"plan item at index "+fmt.Sprint(i)+" contains a blank feature",
				)
			}
		}

		if strings.TrimSpace(item.PaddlePriceID) == "" {
			return domain.NewAppError(
				domain.ErrCodeInvalidParams,
				"paddle price id is required for plan item at index "+fmt.Sprint(i),
			)
		}
	}

	// Plan with variants must have at least 2 items.
	if len(r.Items) > 1 {
		for _, item := range r.Items {
			if len(item.VarOptionIDs) == 0 {
				return domain.NewAppError(domain.ErrCodeInvalidParams, "variant options are required for plan items")
			}
		}
	}

	// Plan without variants can only have one item.
	if len(r.Items) == 1 {
		item := r.Items[0]

		if len(item.VarOptionIDs) > 0 {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "single plan must not contain variant options")
		}
	}
	return nil
}
