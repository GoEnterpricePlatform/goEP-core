package mongo

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *SubscriptionRepository) FindByID(ctx context.Context, id string) (*domain.PaddleSubscription, error) {
	var doc model.PaddleSubscriptionModel
	if err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
		return nil, mapFindError(err, "error finding Paddle subscription")
	}
	return doc.ToDomain(), nil
}

func (r *SubscriptionRepository) FindByPaddleSubscriptionID(ctx context.Context, paddleSubscriptionID string) (*domain.PaddleSubscription, error) {
	var doc model.PaddleSubscriptionModel
	if err := r.Collection.FindOne(ctx, bson.M{"paddle_subscription_id": paddleSubscriptionID}).Decode(&doc); err != nil {
		return nil, mapFindError(err, "error finding Paddle subscription by Paddle ID")
	}
	return doc.ToDomain(), nil
}
