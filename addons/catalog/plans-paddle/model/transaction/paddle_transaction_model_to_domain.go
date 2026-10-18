package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
)

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
