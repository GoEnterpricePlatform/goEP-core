package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *TransactionRepository) CreateIndexes() error {
	_, err := r.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "paddle_transaction_id", Value: 1}},
		Options: options.Index().SetName("paddle_transaction_external_id").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle transaction external ID index: %w", err)
	}
	_, err = r.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "customer_id", Value: 1}, {Key: "updated_at", Value: -1}},
		Options: options.Index().SetName("paddle_transaction_customer_updated_at"),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle transaction customer index: %w", err)
	}
	_, err = r.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "subscription_id", Value: 1}},
		Options: options.Index().SetName("paddle_transaction_subscription_id"),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle transaction subscription index: %w", err)
	}
	return nil
}
