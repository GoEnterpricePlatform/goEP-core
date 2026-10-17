package mongo

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)


func (r *SubscriptionRepository) FindAll(ctx context.Context) ([]*domain.PaddleSubscription, error) {
	return r.find(ctx, bson.M{})
}