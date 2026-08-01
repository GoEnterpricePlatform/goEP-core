package mongo

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/port"
)

var _ port.PaymentProviderRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewPaymentProviderRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: collection,
	}
}
