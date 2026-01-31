package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-tmpl/pkg/features/roles/domain"
	sharedDomain "github.com/amorindev/go-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, role *domain.Role) error {
	id := bson.NewObjectID()
	role.ID = id

	_, err := r.Collection.InsertOne(ctx, role)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting role: %w", sharedDomain.ErrDuplicateKey, err)
		}
		return fmt.Errorf("error inserting role: %w", err)
	}
	role.ID = id.Hex()
	return nil
}
