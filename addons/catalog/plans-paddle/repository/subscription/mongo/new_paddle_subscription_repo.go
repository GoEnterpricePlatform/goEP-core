package mongo

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewPaddleSubscriptionRepo(collection *mongo.Collection) *SubscriptionRepository {
	return &SubscriptionRepository{Collection: collection}
}
