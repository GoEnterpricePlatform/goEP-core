package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"strings"
	"time"
)

func (r *WebhookEventRepository) MarkFailed(ctx context.Context, eventID string, failedAt time.Time, processingErr error) error {
	errorMessage := "unknown processing error"
	if processingErr != nil {
		errorMessage = strings.TrimSpace(processingErr.Error())
	}
	if len(errorMessage) > 2000 {
		errorMessage = errorMessage[:2000]
	}
	result, err := r.Collection.UpdateOne(ctx,
		bson.M{"event_id": eventID, "status": "processing"},
		bson.M{"$set": bson.M{
			"status":     "failed",
			"last_error": errorMessage,
			"failed_at":  failedAt,
		}},
	)
	if err != nil {
		return fmt.Errorf("error marking Paddle webhook event as failed: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("Paddle webhook event %s is not in processing state", eventID)
	}
	return nil
}
