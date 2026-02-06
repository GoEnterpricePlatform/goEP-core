package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-tmpl/pkg/features/permissions/domain"
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

	var permissions []*domain.Permission

	for cursor.Next(ctx) {
		var raw struct {
			ID   bson.ObjectID `bson:"_id"`
			Name string        `bson:"name"`
		}

		if err := cursor.Decode(&raw); err != nil {
			return nil, fmt.Errorf("error decoding permission document: %w", err)
		}

		p := &domain.Permission{
			ID:   raw.ID.Hex(),
			Name: raw.Name,
		}

		permissions = append(permissions, p)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error while iterating permissions: %w", err)
	}

	return permissions, nil
}
