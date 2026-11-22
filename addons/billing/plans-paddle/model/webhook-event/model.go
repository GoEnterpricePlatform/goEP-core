package model

import "time"

type WebhookEventModel struct {
	ID                  string     `bson:"_id"`
	EventID             string     `bson:"event_id"`
	EventType           string     `bson:"event_type"`
	ResourceID          string     `bson:"resource_id,omitempty"`
	OccurredAt          string     `bson:"occurred_at,omitempty"`
	Status              string     `bson:"status"`
	ReceivedAt          time.Time  `bson:"received_at"`
	LastReceivedAt      time.Time  `bson:"last_received_at"`
	ProcessingStartedAt time.Time  `bson:"processing_started_at"`
	ProcessedAt         *time.Time `bson:"processed_at,omitempty"`
	FailedAt            *time.Time `bson:"failed_at,omitempty"`
	Attempts            int        `bson:"attempts"`
	LastError           string     `bson:"last_error,omitempty"`
}
