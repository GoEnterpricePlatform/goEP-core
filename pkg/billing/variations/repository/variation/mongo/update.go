package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) Update(ctx context.Context, variation *domain.Variation) error {
	oID, err := bson.ObjectIDFromHex(variation.ID)
	if err != nil {
		return sharedD.ErrIncorrectID
	}

	update := bson.M{
		"$set": bson.M{
			"name":       variation.Name,
			"updated_at": variation.UpdatedAt,
		},
	}

	filter := bson.M{
		"_id": oID,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var variationModel model.VariationNoSqlModel
	err = r.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&variationModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return sharedD.ErrNotFound
		}
		return fmt.Errorf("failed to update variation: %w", err)
	}

	variationModel.ToDomain(variation)

	return nil
}
