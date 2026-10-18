package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count plans: %w", err)
	}
	return count, nil
}