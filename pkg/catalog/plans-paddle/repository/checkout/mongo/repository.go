package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.PaddleCheckoutRepo = &CheckoutRepository{}

type CheckoutRepository struct {
	Collection *mongo.Collection
}

func NewPaddleCheckoutRepo(collection *mongo.Collection) *CheckoutRepository {
	return &CheckoutRepository{Collection: collection}
}
