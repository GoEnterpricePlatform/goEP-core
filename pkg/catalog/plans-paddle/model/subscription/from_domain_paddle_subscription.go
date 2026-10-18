package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"time"
)

func FromDomainPaddleSubscription(subscription *domain.PaddleSubscription, lastEventAt time.Time) *PaddleSubscriptionModel {
	if subscription == nil {
		return nil
	}
	return &PaddleSubscriptionModel{
		ID:                    subscription.ID,
		PaddleSubscriptionID:  subscription.PaddleSubscriptionID,
		UserID:                subscription.UserID,
		CustomerID:            subscription.CustomerID,
		PlanID:                subscription.PlanID,
		PlanItemID:            subscription.PlanItemID,
		Status:                subscription.Status,
		CurrentPeriodStartsAt: subscription.CurrentPeriodStartsAt,
		CurrentPeriodEndsAt:   subscription.CurrentPeriodEndsAt,
		NextBilledAt:          subscription.NextBilledAt,
		CanceledAt:            subscription.CanceledAt,
		ScheduledChange:       fromDomainSubscriptionChange(subscription.ScheduledChange),
		UpdatedAt:             subscription.UpdatedAt,
		LastEventAt:           lastEventAt,
	}
}
