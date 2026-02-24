package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrIncorrectID
	}

	result, err := r.Collection.DeleteOne(ctx, bson.M{"_id": oID})
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}
