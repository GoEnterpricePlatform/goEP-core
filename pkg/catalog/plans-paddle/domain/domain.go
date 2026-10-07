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
	Order           int                 `json:"order"`
	PaddleProduct   *PaddleProduct      `json:"paddle_product,omitempty"`
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
	PaddlePrice   *PaddlePrice      `json:"paddle_price,omitempty"`
	Options       []*domain.Option  `json:"options"`
	CreatedAt     *time.Time        `json:"created_at"`
	UpdatedAt     *time.Time        `json:"updated_at"`
}

// PaddleProduct contains the Paddle product fields needed to present a plan.
type PaddleProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
	Status      string  `json:"status"`
}

// PaddlePrice contains the Paddle price fields needed to present a plan item.
type PaddlePrice struct {
	ID           string              `json:"id"`
	Name         *string             `json:"name"`
	Description  string              `json:"description"`
	Amount       string              `json:"amount"`
	CurrencyCode string              `json:"currency_code"`
	BillingCycle *PaddleBillingCycle `json:"billing_cycle"`
	Status       string              `json:"status"`
}

type PaddleBillingCycle struct {
	Interval  string `json:"interval"`
	Frequency int    `json:"frequency"`
}
