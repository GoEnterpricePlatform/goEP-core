package model

import "github.com/google/uuid"

// PermissionSQLModel maps the SQL permission row. ID remains provider-neutral
// in the domain and is converted to/from UUID by the PostgreSQL adapter.
type PermissionSQLModel struct {
	ID   uuid.UUID
	Name string
	Type string
}
