package mongo

import (
	"github.com/amorindev/go-cms-tmpl/pkg/posts/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.PostRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewPostRepo(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: collection,
	}
}
