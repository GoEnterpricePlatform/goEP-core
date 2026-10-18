package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *WebhookEventRepository) CreateIndexes() error {
	_, err := r.Collection.Indexes().CreateOne(context.Background(), mongoDriver.IndexModel{
		Keys:    bson.D{{Key: "event_id", Value: 1}},
		Options: options.Index().SetName("paddle_webhook_event_id").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle webhook event ID index: %w", err)
	}
	_, err = r.Collection.Indexes().CreateOne(context.Background(), mongoDriver.IndexModel{
		Keys:    bson.D{{Key: "status", Value: 1}, {Key: "received_at", Value: -1}},
		Options: options.Index().SetName("paddle_webhook_event_status_received_at"),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle webhook event audit index: %w", err)
	}
	return nil
}
