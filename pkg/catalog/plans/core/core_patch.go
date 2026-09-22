package core

import (
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// Note: We use pointers for fields that need to distinguish between
// "not provided" and "explicitly empty".
//
//   - nil → the field was not sent, so it should not be updated.
//   - non-nil → the field was sent and should be updated,
//     even when its value is empty.
//
// This allows partial updates without accidentally overwriting
// fields with zero values.

type PatchPlanReq struct {
	Name        *string              `json:"name"`
	Description *string              `json:"description"`
	PlanItems   *[]*PatchPlanItemReq `json:"plan_items"`
}

type PatchPlanItemReq struct {
	ID           string    `json:"id"`
	Features     *[]string `json:"features"`
	VarOptionIDs *[]string `json:"var_option_ids"`
}

func (req PatchPlanReq) Validate() error {
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)

		if name == "" {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "name is required")
		}
	}

	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)

		if len(desc) > 255 {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "description cannot exceed 255 characters")
		}
	}

	if req.PlanItems != nil {
		if len(*req.PlanItems) == 0 {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "at least one plan item is required")
		}

		for _, item := range *req.PlanItems {
			if item == nil {
				return domain.NewAppError(domain.ErrCodeInvalidParams, "plan item cannot be null")
			}

			if strings.TrimSpace(item.ID) == "" {
				return domain.NewAppError(domain.ErrCodeInvalidParams,"plan item id is required")
			}

			if item.Features != nil {
				for _, feature := range *item.Features {
					if strings.TrimSpace(feature) == "" {
						return domain.NewAppError(domain.ErrCodeInvalidParams,"feature cannot be empty")
					}
				}
			}

			if item.VarOptionIDs != nil {
				for _, ids := range *item.VarOptionIDs {
					if strings.TrimSpace(ids) == "" {
						return domain.NewAppError(domain.ErrCodeInvalidParams,"var_option_id cannot be empty")
					}
				}
			}
		}
	}

	return nil
}
