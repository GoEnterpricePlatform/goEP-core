package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// AssignPermissions adds new permissions to a role without removing existing ones.
// It uses $addToSet to avoid duplicates and preserves previously assigned permissions.
func (r *Repository) AssignPermissions(ctx context.Context, name string, permissionIDs []string) error {
	filter := bson.M{
		"name": name,
	}

	update := bson.M{
		"$addToSet": bson.M{
			"permission_ids": bson.M{
				"$each": permissionIDs,
			},
		},
	}

	result, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(
			"error assigning permissions to role %q: %w",
			name,
			err,
		)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("role %q not found", name)
	}

	return nil
}
