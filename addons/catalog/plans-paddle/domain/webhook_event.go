package domain

// PaddleWebhookEvent is the Catalog-specific envelope for Paddle product events.
type PaddleWebhookEvent struct {
	EventID    string                   `json:"event_id"`
	EventType  string                   `json:"event_type"`
	OccurredAt string                   `json:"occurred_at"`
	Data       PaddleWebhookProductData `json:"data"`
}

type PaddleWebhookProductData struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
