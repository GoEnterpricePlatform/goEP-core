package model

import (
	"time"
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
