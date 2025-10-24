package mongo

import (
	"github.com/amorindev/go-tmpl/pkg/features/users/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Make sure Repository implements ports.PlayerRepository
// at compile time
var _ port.UserRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewUserRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: collection,
	}
}
