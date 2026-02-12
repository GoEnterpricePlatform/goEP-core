package mongo

import (
	port "github.com/amorindev/go-cms-tmpl/pkg/identity/opt-codes/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.OtpCodeRepo = &Repository{}

type Repository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewOtpCodeRepo(client *mongo.Client, otpColl *mongo.Collection) *Repository {
	return &Repository{
		Client:     client,
		Collection: otpColl,
	}
}
