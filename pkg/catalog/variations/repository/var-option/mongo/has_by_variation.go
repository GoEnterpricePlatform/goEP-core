package mongo

import (
	"context"
	"fmt"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) HasByVariation(ctx context.Context, variationID string) (bool, error) {
	id, err := bson.ObjectIDFromHex(variationID)
	if err != nil {
		return false, sharedD.ErrIncorrectID
	}
	err = r.Collection.FindOne(ctx, bson.M{"variation_id": id}).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("checking variation options: %w", err)
	}
	return true, nil
}
