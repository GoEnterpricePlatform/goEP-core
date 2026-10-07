package core

import (
	"fmt"
	"sort"
	"strings"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

type UpdatePaddlePlanReq struct {
	Name            string                  `json:"name"`
	Description     *string                 `json:"description"`
	PaddleProductID string                  `json:"paddle_product_id"`
	Order           int                     `json:"order"`
	Items           []*UpdatePaddlePlanItem `json:"items"`
}

type UpdatePaddlePlanItem struct {
	ID            string   `json:"id"`
	Features      []string `json:"features"`
	VarOptionIDs  []string `json:"var_option_ids"`
	PaddlePriceID string   `json:"paddle_price_id"`
}

func (r UpdatePaddlePlanReq) Validate() error {
	create := CreatePaddlePlanReq{Name: r.Name, Description: r.Description, PaddleProductID: r.PaddleProductID, Order: r.Order}
	for _, item := range r.Items {
		if item == nil {
			create.Items = append(create.Items, nil)
			continue
		}
		create.Items = append(create.Items, &CreatePaddlePlanItem{Features: item.Features, VarOptionIDs: item.VarOptionIDs, PaddlePriceID: item.PaddlePriceID})
	}
	if err := create.Validate(); err != nil {
		return err
	}
	seen := make(map[string]bool, len(r.Items))
	combinations := make(map[string]bool, len(r.Items))
	for i, item := range r.Items {
		optionSeen := make(map[string]bool, len(item.VarOptionIDs))
		for _, id := range item.VarOptionIDs {
			if strings.TrimSpace(id) == "" || optionSeen[id] {
				return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, fmt.Sprintf("invalid variation options at index %d", i))
			}
			optionSeen[id] = true
		}
		optionIDs := append([]string(nil), item.VarOptionIDs...)
		sort.Strings(optionIDs)
		key := strings.Join(optionIDs, ":")
		if combinations[key] {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, fmt.Sprintf("duplicate option combination at index %d", i))
		}
		combinations[key] = true
		if item.ID == "" {
			continue
		}
		if strings.TrimSpace(item.ID) == "" || seen[item.ID] {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, fmt.Sprintf("duplicate or invalid plan item id at index %d", i))
		}
		seen[item.ID] = true
	}
	return nil
}
