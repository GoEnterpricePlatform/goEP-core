package mongo

import (
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.PaddleSubscriptionRepo = &SubscriptionRepository{}

type SubscriptionRepository struct {
	Collection *mongo.Collection
}

func NewPaddleSubscriptionRepo(collection *mongo.Collection) *SubscriptionRepository {
	return &SubscriptionRepository{Collection: collection}
}

func setOptionalString(set bson.M, key string, value *string) {
	if value == nil {
		set[key] = nil
		return
	}
	set[key] = *value
}

func mapFindError(err error, message string) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return sharedD.ErrNotFound
	}
	return fmt.Errorf("%s: %w", message, err)
}
