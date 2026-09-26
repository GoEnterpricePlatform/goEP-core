package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.VariationRepo = &Repository{}

type Repository struct {
	Client         *mongo.Client
	Collection     *mongo.Collection
	VarOptCollName string
}

func NewVariationRepo(client *mongo.Client, collection *mongo.Collection, varOptCollName string) *Repository {
	return &Repository{
		Client:         client,
		Collection:     collection,
		VarOptCollName: varOptCollName,
	}
}
