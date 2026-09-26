package domain

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
)

type PaddlePlan struct {
	ID              string              `json:"id"`
	PlanID          string              `json:"plan_id"`
	Name            string              `json:"name"`
	Description     *string             `json:"description"`
	PaddleProductID string              `json:"paddle_product_id"`
	ImgUrl          *string             `json:"img_url"`
	ImgPath         *string             `json:"-"`
	Items           []*PaddlePlanItem   `json:"items,omitempty"`
	Variations      []*domain.Variation `json:"variations"`
	CreatedAt       *time.Time          `json:"created_at"`
	UpdatedAt       *time.Time          `json:"updated_at"`
}

type PaddlePlanItem struct {
	ID            string            `json:"id"`
	PlanItemID    string            `json:"plan_item_id"`
	ImgUrl        *string           `json:"img_url"`
	ImgPath       *string           `json:"-"`
	Features      []string          `json:"features"`
	Status        domain.PlanStatus `json:"status"`
	VarOptionIDs  []string          `json:"var_option_ids,omitempty"`
	PaddlePriceID string            `json:"paddle_price_id"`
	Options       []*domain.Option  `json:"options"`
	CreatedAt     *time.Time        `json:"created_at"`
	UpdatedAt     *time.Time        `json:"updated_at"`
}
