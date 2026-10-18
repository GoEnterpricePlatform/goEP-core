package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// CreateIndexes creates indexes for Paddle checkouts and subscriptions.
func (r *SubscriptionRepository) CreateIndexes() error {
	_, err := r.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "paddle_subscription_id", Value: 1}},
		Options: options.Index().SetName("paddle_subscription_external_id").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle subscription external ID index: %w", err)
	}
	_, err = r.Collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "status", Value: 1},
				{Key: "updated_at", Value: -1},
			},
			Options: options.Index().
				SetName("paddle_subscription_user_status_updated_at"),
		},
	)
	if err != nil {
		return fmt.Errorf("error creating Paddle subscription user index: %w", err)
	}
	_, err = r.Collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "customer_id", Value: 1},
				{Key: "status", Value: 1},
				{Key: "updated_at", Value: -1},
			},
			Options: options.Index().SetName("paddle_subscription_customer_status_updated_at"),
		},
	)
	if err != nil {
		return fmt.Errorf("error creating Paddle subscription customer index: %w", err)
	}

	return nil
}
