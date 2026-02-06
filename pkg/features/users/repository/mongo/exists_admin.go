package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-tmpl/pkg/features/roles/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) ExistsAdmin(ctx context.Context) (bool, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "roles"},
				{Key: "localField", Value: "role_ids"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "roles"},
			}},
		},
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "roles.name", Value: string(domain.RoleSystemAdmin)},
			}},
		},
		bson.D{
			{Key: "$limit", Value: 1},
		},
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return false, fmt.Errorf("failed to execute aggregate pipeline: %w", err)
	}

	if cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return false, fmt.Errorf("failed to decode aggregate result: %w", err)
		}
		return true, nil
	}

	return false, nil
}
