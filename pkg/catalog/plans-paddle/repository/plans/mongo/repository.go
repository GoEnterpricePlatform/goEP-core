package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.PaddlePlanRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewPlanPaddleRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: collection,
	}
}
