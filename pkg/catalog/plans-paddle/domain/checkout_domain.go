package domain

type PaddleCheckout struct {
	TransactionID  string  `json:"transaction_id"`
	PlanID         string  `json:"plan_id"`
	PlanItemID     string  `json:"plan_item_id"`
	PaddlePriceID  string  `json:"paddle_price_id"`
	Status         string  `json:"status"`
	SubscriptionID *string `json:"subscription_id,omitempty"`
	UpdatedAt      string  `json:"updated_at"`
}
