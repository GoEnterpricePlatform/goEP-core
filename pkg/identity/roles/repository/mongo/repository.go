package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.RoleRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewRoleRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: collection,
	}
}
