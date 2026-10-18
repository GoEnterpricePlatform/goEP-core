package helpers

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
)

func PreserveTransactionFields(incoming, current *domain.PaddleTransaction) {
	if incoming.ID == "" {
		incoming.ID = current.ID
	}
	if incoming.PaddleTransactionID == "" {
		incoming.PaddleTransactionID = current.PaddleTransactionID
	}
	if incoming.UserID == "" {
		incoming.UserID = current.UserID
	}
	if incoming.CustomerID == "" {
		incoming.CustomerID = current.CustomerID
	}
	if incoming.PlanID == "" {
		incoming.PlanID = current.PlanID
	}
	if incoming.PlanItemID == "" {
		incoming.PlanItemID = current.PlanItemID
	}
	if incoming.SubscriptionID == nil {
		incoming.SubscriptionID = current.SubscriptionID
	}
	if current.CurrencyCode != "" {
		incoming.CurrencyCode = current.CurrencyCode
	}
	if current.Subtotal != 0 {
		incoming.Subtotal = current.Subtotal
	}
	if current.Tax != 0 {
		incoming.Tax = current.Tax
	}
	if current.Total != 0 {
		incoming.Total = current.Total
	}
	if current.PaymentDetailsProvided {
		incoming.PaymentDetailsProvided = true
	}
}
