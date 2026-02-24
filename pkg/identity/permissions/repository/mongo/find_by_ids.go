package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/permissions/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/permissions/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// FindByIDs returns the permissions for the given permission IDs.
//
// This function is typically used after resolving roles that contain a list
// of permission IDs. It queries the database for those permission documents
// and returns the full Permission domain entities. If no IDs are provided,
// it returns an empty slice.
func (r *Repository) FindByIDs(ctx context.Context, permissionIDs []string) ([]*domain.Permission, error) {
	if len(permissionIDs) == 0 {
		return []*domain.Permission{}, nil
	}

	var permissionOIDs []bson.ObjectID
	for _, id := range permissionIDs {
		oID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid permission ID format (%s): %w", id, err)
		}
		permissionOIDs = append(permissionOIDs, oID)
	}

	filter := bson.M{"_id": bson.M{"$in": permissionOIDs}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting permissions by IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var permissions []*domain.Permission
	for cursor.Next(ctx) {
		var permissionModel model.PermissionNoSqlModel
		if err := cursor.Decode(&permissionModel); err != nil {
			return nil, fmt.Errorf("error decoding permission document: %w", err)
		}

		var permission domain.Permission
		permissionModel.ToDomain(&permission)

		permissions = append(permissions, &permission)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error while iterating permissions documents: %w", err)
	}

	return permissions, nil
}
