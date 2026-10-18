package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

func (r *WebhookEventRepository) MarkProcessed(ctx context.Context, eventID string, processedAt time.Time) error {
	result, err := r.Collection.UpdateOne(ctx,
		bson.M{"event_id": eventID, "status": "processing"},
		bson.M{
			"$set":   bson.M{"status": "processed", "processed_at": processedAt},
			"$unset": bson.M{"last_error": "", "failed_at": ""},
		},
	)
	if err != nil {
		return fmt.Errorf("error marking Paddle webhook event as processed: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("Paddle webhook event %s is not in processing state", eventID)
	}
	return nil
}
