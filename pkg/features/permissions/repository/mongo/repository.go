package mongo

import (
	"github.com/amorindev/go-tmpl/pkg/features/permissions/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.PermissionRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewPermissionRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
    return &Repository{
		Client:       client,
		Collection:   collection,
	}
}
