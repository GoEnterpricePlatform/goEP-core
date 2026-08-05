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

func (r *Repository) Update(ctx context.Context, varOption *domain.VarOption) error {
	oID, err := bson.ObjectIDFromHex(varOption.ID)
	if err != nil {
		return sharedD.ErrIncorrectID
	}

	variationOID, err := bson.ObjectIDFromHex(varOption.VariationID)
	if err != nil {
		return sharedD.ErrIncorrectID
	}

	update := bson.M{
		"$set": bson.M{
			"label":      varOption.Label,
			"value":      varOption.Value,
			"updated_at": varOption.UpdatedAt,
		},
	}

	filter := bson.M{
		"_id":          oID,
		"variation_id": variationOID,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var varOptionModel model.VarOptionNoSqlModel
	err = r.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&varOptionModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return sharedD.ErrNotFound
		}
		return fmt.Errorf("failed to update varOption: %w", err)
	}

	varOptionModel.ToDomain(varOption)

	return nil
}
