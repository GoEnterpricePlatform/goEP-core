package mongo

import (
	"context"
	"fmt"

	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%w: invalid otpID :%w", sharedD.ErrIncorrectID, err)
	}

	result, err := r.Collection.DeleteOne(ctx, bson.M{"_id": oID})
	if err != nil {
		return fmt.Errorf("error deleting otpCode: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%w: otpCode not found", sharedD.ErrNotFound)
	}
	return nil
}
