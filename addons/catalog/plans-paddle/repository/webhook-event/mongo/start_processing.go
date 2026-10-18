package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/model/webhook-event"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
)

// StartProcessing claims a new event, retries a failed/stale event, or reports
// that an already processed event can safely be acknowledged without replaying it.
func (r *WebhookEventRepository) StartProcessing(ctx context.Context, event *domain.PaddleWebhookEvent, now time.Time) (bool, error) {
	if event == nil || event.EventID == "" {
		return false, fmt.Errorf("Paddle webhook event is missing its event ID")
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}

	record := model.WebhookEventModel{
		ID:                  bson.NewObjectID().Hex(),
		EventID:             event.EventID,
		EventType:           event.EventType,
		ResourceID:          event.Data.ID,
		OccurredAt:          event.OccurredAt,
		Status:              "processing",
		ReceivedAt:          now,
		LastReceivedAt:      now,
		ProcessingStartedAt: now,
		Attempts:            1,
	}
	if _, err := r.Collection.InsertOne(ctx, record); err == nil {
		return false, nil
	} else if !mongoDriver.IsDuplicateKeyError(err) {
		return false, fmt.Errorf("error inserting Paddle webhook event: %w", err)
	}

	// A duplicate delivery may retry an earlier failure. Also reclaim processing
	// records left behind by a crashed worker after the lease expires.
	result, err := r.Collection.UpdateOne(ctx,
		bson.M{
			"event_id": event.EventID,
			"$or": bson.A{
				bson.M{"status": "failed"},
				bson.M{
					"status":                "processing",
					"processing_started_at": bson.M{"$lte": now.Add(-processingLease)},
				},
			},
		},
		bson.M{
			"$set": bson.M{
				"event_type":            event.EventType,
				"resource_id":           event.Data.ID,
				"occurred_at":           event.OccurredAt,
				"status":                "processing",
				"last_received_at":      now,
				"processing_started_at": now,
			},
			"$inc":   bson.M{"attempts": 1},
			"$unset": bson.M{"last_error": "", "processed_at": "", "failed_at": ""},
		},
	)
	if err != nil {
		return false, fmt.Errorf("error claiming Paddle webhook event retry: %w", err)
	}
	if result.ModifiedCount > 0 {
		return false, nil
	}

	var existing model.WebhookEventModel
	if err := r.Collection.FindOne(ctx, bson.M{"event_id": event.EventID}).Decode(&existing); err != nil {
		return false, fmt.Errorf("error checking duplicate Paddle webhook event: %w", err)
	}
	if existing.Status == "processed" {
		return true, nil
	}
	return false, fmt.Errorf("Paddle webhook event %s is already being processed", event.EventID)
}
