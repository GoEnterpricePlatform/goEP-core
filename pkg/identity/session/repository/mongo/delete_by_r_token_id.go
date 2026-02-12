package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) DeleteByRTokenID(ctx context.Context, rTokenID string) error {
	filter := bson.M{"refresh_token_id": rTokenID}

	result, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

    if result.DeletedCount == 0 {
        return domain.ErrNotFound
    }

	return nil
}
