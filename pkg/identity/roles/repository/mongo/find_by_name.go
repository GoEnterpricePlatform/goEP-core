package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/roles/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/roles/model"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	var userModel model.RoleNoSqlModel

	filter := bson.D{{Key: "name", Value: name}}

	err := r.Collection.FindOne(ctx, filter).Decode(&userModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: role with name %s not found: %w", sharedD.ErrNotFound, name, err)
		}
		return nil, fmt.Errorf("role with name %s not found: %w", name, err)
	}

	
	var role domain.Role
	userModel.ToDomain(&role)

	return &role, nil
}
