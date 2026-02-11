package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/features/roles/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role

	filter := bson.D{{Key: "name", Value: name}}

	err := r.Collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: role with name %s not found: %w", sharedD.ErrNotFound, name, err)
		}
		return nil, fmt.Errorf("role with name %s not found: %w", name, err)
	}

	oID, ok := role.ID.(bson.ObjectID)
	if !ok {
		return nil, fmt.Errorf("role ID is not a bson.ObjectID")
	}
	role.ID = oID.Hex()

	return &role, nil
}
