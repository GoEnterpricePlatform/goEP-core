package mongo

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model/subscription"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *SubscriptionRepository) FindByID(ctx context.Context, id string) (*domain.PaddleSubscription, error) {
	var doc model.PaddleSubscriptionModel
	if err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
		return nil, mapFindError(err, "error finding Paddle subscription")
	}
	return doc.ToDomain(), nil
}
