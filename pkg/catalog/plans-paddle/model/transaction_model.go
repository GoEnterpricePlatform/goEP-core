package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
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

func (m *PaddleTransactionModel) ToDomain() *domain.PaddleTransaction {
	if m == nil {
		return nil
	}
	return &domain.PaddleTransaction{
		ID:                     m.ID,
		PaddleTransactionID:    m.PaddleTransactionID,
		UserID:                 m.UserID,
		SubscriptionID:         m.SubscriptionID,
		CustomerID:             m.CustomerID,
		PlanID:                 m.PlanID,
		PlanItemID:             m.PlanItemID,
		Status:                 m.Status,
		CurrencyCode:           m.CurrencyCode,
		Subtotal:               m.Subtotal,
		Tax:                    m.Tax,
		Total:                  m.Total,
		UpdatedAt:              m.UpdatedAt,
		LastEventAt:            m.LastEventAt,
		PaymentDetailsProvided: m.PaymentDetailsProvided,
	}
}

func FromDomainPaddleTransaction(transaction *domain.PaddleTransaction) *PaddleTransactionModel {
	if transaction == nil {
		return nil
	}
	return &PaddleTransactionModel{
		ID:                     transaction.ID,
		PaddleTransactionID:    transaction.PaddleTransactionID,
		UserID:                 transaction.UserID,
		SubscriptionID:         transaction.SubscriptionID,
		CustomerID:             transaction.CustomerID,
		PlanID:                 transaction.PlanID,
		PlanItemID:             transaction.PlanItemID,
		Status:                 transaction.Status,
		CurrencyCode:           transaction.CurrencyCode,
		Subtotal:               transaction.Subtotal,
		Tax:                    transaction.Tax,
		Total:                  transaction.Total,
		UpdatedAt:              transaction.UpdatedAt,
		LastEventAt:            transaction.LastEventAt,
		PaymentDetailsProvided: transaction.PaymentDetailsProvided,
	}
}
