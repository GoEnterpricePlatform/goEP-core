package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *CustomerRepository) CreateIndexes() error {
	_, err := r.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetName("paddle_customer_user_id").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle customer user index: %w", err)
	}
	_, err = r.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "paddle_customer_id", Value: 1}},
		Options: options.Index().SetName("paddle_customer_external_id").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle customer external ID index: %w", err)
	}
	return nil
}
