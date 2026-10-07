package mongo

import (
	"context"
	"fmt"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) IsOptionInUse(ctx context.Context, optionID string) (bool, error) {
	id, err := bson.ObjectIDFromHex(optionID)
	if err != nil {
		return false, sharedD.ErrIncorrectID
	}
	err = r.Collection.FindOne(ctx, bson.M{"items.var_option_ids": id}).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("checking plan option references: %w", err)
	}
	return true, nil
}
