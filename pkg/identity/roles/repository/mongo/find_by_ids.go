package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// FindByIDs returns the roles for the given role IDs.
//
// This function is used during the sign-in process: after retrieving a user
// that contains a list of role IDs, it queries the database for those role
// documents and returns the full Role entities. If no IDs are provided,
// it returns an empty slice.
func (r *Repository) FindByIDs(ctx context.Context, roleIDs []string) ([]*domain.Role, error) {
	if len(roleIDs) == 0 {
		return []*domain.Role{}, nil
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
		return nil, fmt.Errorf("error getting roles by IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var roles []*domain.Role
	for cursor.Next(ctx) {
		var roleModel model.RoleNoSqlModel
		if err := cursor.Decode(&roleModel); err != nil {
			return nil, fmt.Errorf("error decoding role document: %w", err)
		}

		var role domain.Role
		roleModel.ToDomain(&role)

		roles = append(roles, &role)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error while iterating roles documents: %w", err)
	}

	return roles, nil
}
