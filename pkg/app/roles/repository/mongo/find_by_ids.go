package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-tmpl/pkg/app/roles/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// FindByIDs returns the names of roles for the given role IDs.
//
// This function is used during the sign-in process: after retrieving a user
// that contains a list of role IDs, it queries the database for those role
// documents and returns their names. If no IDs are provided, it returns an
// empty slice.
func (r *Repository) FindByIDs(ctx context.Context, roleIDs []string) ([]string, error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	var roleOIDs []bson.ObjectID
	for _, id := range roleIDs {
		oID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid role ID format (%s): %w", id, err)
		}
		roleOIDs = append(roleOIDs, oID)
	}

	filter := bson.M{"_id": bson.M{"$in": roleOIDs}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles by IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var roleNames []string
	for cursor.Next(ctx) {
		var role domain.Role
		if err := cursor.Decode(&role); err != nil {
			return nil, fmt.Errorf("failed to decode role document: %w", err)
		}
		roleNames = append(roleNames, role.Name)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error while iterating role documents: %w", err)
	}

	return roleNames, nil
}
