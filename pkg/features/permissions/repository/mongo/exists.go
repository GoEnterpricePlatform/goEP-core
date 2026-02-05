package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Exists checks if a permission with the given name already exists in the collection.
func (r *Repository) Exists(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}

	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("error checking permission existence: %w", err)
	}
	return count > 0, nil
}