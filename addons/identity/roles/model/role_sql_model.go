package model

import "github.com/google/uuid"

// RoleSQLModel maps tb_identity_role. Its permission IDs are stored in the
// separate tb_identity_role_permission relation.
type RoleSQLModel struct {
	ID   uuid.UUID
	Name string
}
type RolePermissionSQLModel struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}
