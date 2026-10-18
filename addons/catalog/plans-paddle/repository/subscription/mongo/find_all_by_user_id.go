package mongo

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *SubscriptionRepository) FindAllByUserID(ctx context.Context, userID string) ([]*domain.PaddleSubscription, error) {
	return r.find(ctx, bson.M{"user_id": userID})
}
