package mongo

import (
	"github.com/amorindev/go-tmpl/pkg/features/session/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.SessionRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewSessionRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: collection,
	}
}
