package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Exists checks if a user with the given id already exists in the collection.
func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("%w: failed to convert userID to ObjectID :%w", domain.ErrIncorrectID, err)
	}

	filter := bson.M{"_id": oID}

	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %s", err.Error())
	}

	return count > 0, nil
}
