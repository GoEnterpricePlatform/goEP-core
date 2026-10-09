package model

import "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"

type PaddleCheckoutModel struct {
	TransactionID  string  `bson:"_id"`
	PlanID         string  `bson:"plan_id"`
	PlanItemID     string  `bson:"plan_item_id"`
	PaddlePriceID  string  `bson:"paddle_price_id"`
	Status         string  `bson:"status"`
	SubscriptionID *string `bson:"subscription_id,omitempty"`
	UpdatedAt      string  `bson:"updated_at"`
}

func (m *PaddleCheckoutModel) ToDomain(checkout *domain.PaddleCheckout) {
	if m == nil || checkout == nil {
		return
	}

	checkout.TransactionID = m.TransactionID
	checkout.PlanID = m.PlanID
	checkout.PlanItemID = m.PlanItemID
	checkout.PaddlePriceID = m.PaddlePriceID
	checkout.Status = m.Status
	checkout.SubscriptionID = m.SubscriptionID
	checkout.UpdatedAt = m.UpdatedAt
}

func FromDomainPaddleCheckout(checkout *domain.PaddleCheckout) *PaddleCheckoutModel {
	if checkout == nil {
		return nil
	}

	return &PaddleCheckoutModel{
		TransactionID:  checkout.TransactionID,
		PlanID:         checkout.PlanID,
		PlanItemID:     checkout.PlanItemID,
		PaddlePriceID:  checkout.PaddlePriceID,
		Status:         checkout.Status,
		SubscriptionID: checkout.SubscriptionID,
		UpdatedAt:      checkout.UpdatedAt,
	}
}
