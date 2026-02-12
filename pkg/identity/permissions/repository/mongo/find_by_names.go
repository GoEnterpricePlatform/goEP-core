package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/permissions/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/permissions/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) FindByNames(ctx context.Context, names []domain.PermissionName) ([]*domain.Permission, error) {
	if len(names) == 0 {
		return []*domain.Permission{}, nil
	}

	filter := bson.M{
		"name": bson.M{
			"$in": names,
		},
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting permissions: %w", err)
	}
	defer cursor.Close(ctx)

	var permissionsModel []*model.PermissionNoSqlModel

	for cursor.Next(ctx) {
		var permssionModel model.PermissionNoSqlModel

		if err := cursor.Decode(&permssionModel); err != nil {
			return nil, fmt.Errorf("error decoding permission document: %w", err)
		}

		permissionsModel = append(permissionsModel, &permssionModel)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error while iterating permissions: %w", err)
	}

	return model.ToDomainList(permissionsModel), nil
}
