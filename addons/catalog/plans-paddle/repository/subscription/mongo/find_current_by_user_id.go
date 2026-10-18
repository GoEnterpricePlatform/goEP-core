package mongo

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/model/subscription"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *SubscriptionRepository) FindCurrentByUserID(ctx context.Context, userID string) (*domain.PaddleSubscription, error) {
	var doc model.PaddleSubscriptionModel
	opts := options.FindOne().SetSort(bson.D{{Key: "updated_at", Value: -1}, {Key: "_id", Value: -1}})
	err := r.Collection.FindOne(ctx, bson.M{
		"user_id": userID,
		"status":  bson.M{"$in": bson.A{"active", "trialing"}},
	}, opts).Decode(&doc)
	if err != nil {
		return nil, mapFindError(err, "error finding current Paddle subscription")
	}
	return doc.ToDomain(), nil
}
