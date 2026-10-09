package domain

type PaddleWebhookEvent struct {
	EventType string                 `json:"event_type"`
	Data      PaddleWebhookEventData `json:"data"`
}

type PaddleWebhookEventData struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	SubscriptionID *string `json:"subscription_id"`
}
