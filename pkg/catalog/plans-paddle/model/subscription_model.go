package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
)

type PaddleSubscriptionModel struct {
	ID                    string                         `bson:"_id"`
	PaddleSubscriptionID  string                         `bson:"paddle_subscription_id"`
	UserID                string                         `bson:"user_id"`
	CustomerID            string                         `bson:"customer_id,omitempty"`
	PlanID                string                         `bson:"plan_id,omitempty"`
	PlanItemID            string                         `bson:"plan_item_id,omitempty"`
	Status                string                         `bson:"status"`
	CurrentPeriodStartsAt *string                        `bson:"current_period_starts_at,omitempty"`
	CurrentPeriodEndsAt   *string                        `bson:"current_period_ends_at,omitempty"`
	NextBilledAt          *string                        `bson:"next_billed_at,omitempty"`
	CanceledAt            *string                        `bson:"canceled_at,omitempty"`
	ScheduledChange       *PaddleSubscriptionChangeModel `bson:"scheduled_change,omitempty"`
	UpdatedAt             string                         `bson:"updated_at"`
	LastEventAt           time.Time                      `bson:"last_event_at,omitempty"`
}

type PaddleSubscriptionChangeModel struct {
	Action      string `bson:"action"`
	EffectiveAt string `bson:"effective_at"`
}

func (m *PaddleSubscriptionModel) ToDomain() *domain.PaddleSubscription {
	if m == nil {
		return nil
	}
	return &domain.PaddleSubscription{
		ID:                    m.ID,
		PaddleSubscriptionID:  m.PaddleSubscriptionID,
		UserID:                m.UserID,
		CustomerID:            m.CustomerID,
		PlanID:                m.PlanID,
		PlanItemID:            m.PlanItemID,
		Status:                m.Status,
		CurrentPeriodStartsAt: m.CurrentPeriodStartsAt,
		CurrentPeriodEndsAt:   m.CurrentPeriodEndsAt,
		NextBilledAt:          m.NextBilledAt,
		CanceledAt:            m.CanceledAt,
		ScheduledChange:       toDomainSubscriptionChange(m.ScheduledChange),
		UpdatedAt:             m.UpdatedAt,
	}
}

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

func toDomainSubscriptionChange(change *PaddleSubscriptionChangeModel) *domain.PaddleSubscriptionChange {
	if change == nil {
		return nil
	}
	return &domain.PaddleSubscriptionChange{Action: change.Action, EffectiveAt: change.EffectiveAt}
}

func fromDomainSubscriptionChange(change *domain.PaddleSubscriptionChange) *PaddleSubscriptionChangeModel {
	if change == nil {
		return nil
	}
	return &PaddleSubscriptionChangeModel{Action: change.Action, EffectiveAt: change.EffectiveAt}
}
