package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) FindAll(ctx context.Context) ([]*domain.Role, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error getting roles: %w", err)
	}
	defer cursor.Close(ctx)

	var roles []*domain.Role
	for cursor.Next(ctx) {
		var roleModel model.RoleNoSqlModel
		err := cursor.Decode(&roleModel)
		if err != nil {
			return nil, fmt.Errorf("error decoding role: %w", err)
		}

		var role domain.Role
		roleModel.ToDomain(&role)

		roles = append(roles, &role)
	}

	return roles, nil
}
