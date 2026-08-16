package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.VarOptionRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewVarOptionRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: collection,
	}
}
