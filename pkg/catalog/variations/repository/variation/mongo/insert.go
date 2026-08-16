package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, variation *domain.Variation) error {
	id := bson.NewObjectID()

	variationModel, err := model.FromDomain(variation, id)
	if err != nil {
		return fmt.Errorf("error mapping variation to mongo model: %w", err)
	}

	_, err = r.Collection.InsertOne(ctx, variationModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting variation: %s", sharedD.ErrDuplicateKey, err.Error())
		}
		return fmt.Errorf("error inserting variation: %s", err.Error())
	}

	variationModel.ToDomain(variation)

	return nil
}
