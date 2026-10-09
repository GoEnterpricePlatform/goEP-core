package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *CheckoutRepository) UpdateStatus(ctx context.Context, transactionID, status string, subscriptionID *string, updatedAt string) error {
	set := bson.M{"status": status, "updated_at": updatedAt}
	if subscriptionID != nil {
		set["subscription_id"] = *subscriptionID
	}
	result, err := r.Collection.UpdateOne(ctx, bson.M{
		"_id":    transactionID,
		"status": bson.M{"$ne": "completed"},
	}, bson.M{"$set": set})
	if err != nil {
		return fmt.Errorf("error updating Paddle checkout status: %w", err)
	}
	if result.MatchedCount == 0 {
		// Duplicate notifications and events for unknown transactions are safe to acknowledge.
		return nil
	}
	return nil
}
