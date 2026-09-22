package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.PlanRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewPlanRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: collection,
	}
}
