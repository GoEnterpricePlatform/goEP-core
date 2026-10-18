package domain

import "time"

type PaddleTransaction struct {
	ID                     string    `json:"id"`
	PaddleTransactionID    string    `json:"paddle_transaction_id"`
	UserID                 string    `json:"-"`
	SubscriptionID         *string   `json:"subscription_id,omitempty"`
	CustomerID             string    `json:"customer_id,omitempty"`
	PlanID                 string    `json:"plan_id"`
	PlanItemID             string    `json:"plan_item_id"`
	Status                 string    `json:"status"`
	CurrencyCode           string    `json:"currency_code"`
	Subtotal               int64     `json:"subtotal"`
	Tax                    int64     `json:"tax"`
	Total                  int64     `json:"total"`
	UpdatedAt              string    `json:"updated_at"`
	LastEventAt            time.Time `json:"-"`
	PaymentDetailsProvided bool      `json:"-"`
}
