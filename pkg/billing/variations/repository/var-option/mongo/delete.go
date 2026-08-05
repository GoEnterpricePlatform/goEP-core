package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) Delete(ctx context.Context, id string, variationID string) error {
    oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrIncorrectID
	}

    variationOID, err := bson.ObjectIDFromHex(variationID)
	if err != nil {
		return domain.ErrIncorrectID
	}

    filter := bson.M{
        "_id":oID,
        "variation_id": variationOID,
    }

    result, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting varOption: %w", err)
	}

    if result.DeletedCount == 0 {
        return domain.ErrNotFound
    }

    return nil
}
