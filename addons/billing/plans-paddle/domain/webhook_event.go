package domain

type PaddleWebhookEvent struct {
	EventID    string                 `json:"event_id"`
	EventType  string                 `json:"event_type"`
	OccurredAt string                 `json:"occurred_at"`
	Data       PaddleWebhookEventData `json:"data"`
}

type PaddleWebhookEventData struct {
	ID                   string                           `json:"id"`
	CustomerID           string                           `json:"customer_id"`
	CurrencyCode         string                           `json:"currency_code"`
	Status               string                           `json:"status"`
	UpdatedAt            string                           `json:"updated_at"`
	SubscriptionID       *string                          `json:"subscription_id"`
	CustomData           map[string]any                   `json:"custom_data"`
	CurrentBillingPeriod *PaddleWebhookTimePeriod         `json:"current_billing_period"`
	NextBilledAt         *string                          `json:"next_billed_at"`
	CanceledAt           *string                          `json:"canceled_at"`
	ScheduledChange      *PaddleWebhookScheduledChange    `json:"scheduled_change"`
	Items                []PaddleWebhookItem              `json:"items"`
	Details              *PaddleWebhookTransactionDetails `json:"details"`
}

type PaddleWebhookTransactionDetails struct {
	Totals PaddleWebhookTransactionTotals `json:"totals"`
}

type PaddleWebhookTransactionTotals struct {
	Subtotal string `json:"subtotal"`
	Tax      string `json:"tax"`
	Total    string `json:"total"`
}

type PaddleWebhookTimePeriod struct {
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type PaddleWebhookScheduledChange struct {
	Action      string `json:"action"`
	EffectiveAt string `json:"effective_at"`
}

type PaddleWebhookItem struct {
	Price PaddleWebhookPrice `json:"price"`
}

type PaddleWebhookPrice struct {
	ID string `json:"id"`
}
