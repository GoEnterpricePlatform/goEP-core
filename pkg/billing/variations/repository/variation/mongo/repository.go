package mongo

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/port"
)

var _ port.VariationRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewVariationRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: collection,
	}
}