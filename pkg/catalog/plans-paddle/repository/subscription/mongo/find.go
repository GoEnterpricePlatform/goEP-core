package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model/subscription"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *SubscriptionRepository) find(ctx context.Context, filter bson.M) ([]*domain.PaddleSubscription, error) {
	cursor, err := r.Collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("error finding Paddle subscriptions: %w", err)
	}
	defer cursor.Close(ctx)
	var docs []model.PaddleSubscriptionModel
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("error decoding Paddle subscriptions: %w", err)
	}
	subscriptions := make([]*domain.PaddleSubscription, 0, len(docs))
	for i := range docs {
		subscriptions = append(subscriptions, docs[i].ToDomain())
	}
	return subscriptions, nil
}
