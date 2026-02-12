package domain

type Role struct {
	ID            string   `json:"id" bson:"_id"`
	Name          string        `json:"name" bson:"name"`
	PermissionIDs []string `bson:"permission_ids"`
}

// Initialize PermissionIDs as an empty slice to allow MongoDB $addToSet operations.
func NewRole(name string) *Role {
	return &Role{
		Name:          name,
		PermissionIDs: []string{},
	}
}

// system roles should be saved in the db using initializer
type RoleName string

const (
	RoleSystemAdmin RoleName = "ROLE_SYSTEM_ADMIN"
	RoleUser        RoleName = "USER"
)
