package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) CreateIndexes() error {
	_, err := r.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "desc", Value: "text"},
		},
		Options: options.Index().SetName("post_text_search").SetWeights(bson.D{
			{Key: "title", Value: 10},
			{Key: "desc", Value: 5},
		}),
	})
	if err != nil {
		return fmt.Errorf("error creating post text index: %w", err)
	}
	return nil
}
