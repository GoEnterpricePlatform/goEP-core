package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoleNoSqlModel struct {
	ID            bson.ObjectID   `bson:"_id"`
	Name          string          `bson:"name"`
	PermissionIDs []bson.ObjectID `bson:"permission_ids"`
}

func (m *RoleNoSqlModel) ToDomain(r *domain.Role) {
	if m == nil {
		return
	}

	permissionIDs := make([]string, 0, len(m.PermissionIDs))
	for _, p := range m.PermissionIDs {
		permissionIDs = append(permissionIDs, p.Hex())
	}

	r.ID = m.ID.Hex()
	r.Name = m.Name
	r.PermissionIDs = permissionIDs
}

func FromDomain(r *domain.Role, id bson.ObjectID) (*RoleNoSqlModel, error) {
	if r == nil {
		return nil, nil
	}

	permissionObjectIDs := make([]bson.ObjectID, 0, len(r.PermissionIDs))

	for _, p := range r.PermissionIDs {
		objID, err := bson.ObjectIDFromHex(p)
		if err != nil {
			return nil, err
		}
		permissionObjectIDs = append(permissionObjectIDs, objID)
	}

	return &RoleNoSqlModel{
		ID:            id,
		Name:          r.Name,
		PermissionIDs: permissionObjectIDs,
	}, nil
}
