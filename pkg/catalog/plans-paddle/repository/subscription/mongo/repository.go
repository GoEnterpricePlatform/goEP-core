package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.PaddleSubscriptionRepo = &SubscriptionRepository{}

type SubscriptionRepository struct {
	Collection *mongo.Collection
}
