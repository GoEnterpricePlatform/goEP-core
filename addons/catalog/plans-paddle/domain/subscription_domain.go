package domain

type PaddleSubscription struct {
	ID                    string                    `json:"id"`
	PaddleSubscriptionID  string                    `json:"paddle_subscription_id"`
	UserID                string                    `json:"user_id"`
	CustomerID            string                    `json:"customer_id,omitempty"`
	PlanID                string                    `json:"plan_id,omitempty"`
	PlanItemID            string                    `json:"plan_item_id,omitempty"`
	Status                string                    `json:"status"`
	CurrentPeriodStartsAt *string                   `json:"current_period_starts_at,omitempty"`
	CurrentPeriodEndsAt   *string                   `json:"current_period_ends_at,omitempty"`
	NextBilledAt          *string                   `json:"next_billed_at,omitempty"`
	CanceledAt            *string                   `json:"canceled_at,omitempty"`
	ScheduledChange       *PaddleSubscriptionChange `json:"scheduled_change,omitempty"`
	UpdatedAt             string                    `json:"updated_at"`
}

type PaddleSubscriptionChange struct {
	Action      string `json:"action"`
	EffectiveAt string `json:"effective_at"`
}

type PaddleSubscriptionAccess struct {
	HasAccess    bool                `json:"has_access"`
	Subscription *PaddleSubscription `json:"subscription,omitempty"`
}
