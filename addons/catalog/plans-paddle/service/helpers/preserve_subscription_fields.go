package helpers

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
)

func PreserveSubscriptionFields(incoming, current *domain.PaddleSubscription) {
	if incoming.ID == "" {
		incoming.ID = current.ID
	}
	if incoming.PaddleSubscriptionID == "" {
		incoming.PaddleSubscriptionID = current.PaddleSubscriptionID
	}
	if incoming.UserID == "" {
		incoming.UserID = current.UserID
	}
	if incoming.PlanID == "" {
		incoming.PlanID = current.PlanID
	}
	if incoming.PlanItemID == "" {
		incoming.PlanItemID = current.PlanItemID
	}
	if incoming.CustomerID == "" {
		incoming.CustomerID = current.CustomerID
	}
}
