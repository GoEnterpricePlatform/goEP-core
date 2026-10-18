package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
)

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
