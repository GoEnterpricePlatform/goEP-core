package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, fmt.Errorf("error counting paddle plans: %w", err)
	}
	return count, nil
}
