package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.PlanRepo = &Repository{}

type Repository struct {
	Client             *mongo.Client
	Collection         *mongo.Collection
	VarOptCollName     string
	variationsCollName string
}

func NewPlanRepo(client *mongo.Client, collection *mongo.Collection, varOptCollName string, variationsCollName string) *Repository {
	return &Repository{
		Client:             client,
		Collection:         collection,
		VarOptCollName:     varOptCollName,
		variationsCollName: variationsCollName,
	}
}
