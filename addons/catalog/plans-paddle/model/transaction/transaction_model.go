package model

import (
	"time"
)

type PaddleTransactionModel struct {
	ID                     string    `bson:"_id"`
	PaddleTransactionID    string    `bson:"paddle_transaction_id"`
	UserID                 string    `bson:"user_id,omitempty"`
	SubscriptionID         *string   `bson:"subscription_id,omitempty"`
	CustomerID             string    `bson:"customer_id,omitempty"`
	PlanID                 string    `bson:"plan_id"`
	PlanItemID             string    `bson:"plan_item_id"`
	Status                 string    `bson:"status"`
	CurrencyCode           string    `bson:"currency_code,omitempty"`
	Subtotal               int64     `bson:"subtotal,omitempty"`
	Tax                    int64     `bson:"tax,omitempty"`
	Total                  int64     `bson:"total,omitempty"`
	UpdatedAt              string    `bson:"updated_at"`
	LastEventAt            time.Time `bson:"last_event_at,omitempty"`
	PaymentDetailsProvided bool      `bson:"payment_details_provided,omitempty"`
}
