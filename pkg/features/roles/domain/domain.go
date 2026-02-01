package domain

type Role struct {
	ID   interface{} `json:"id" bson:"_id"`
	Name string      `json:"name" bson:"name"`
}

func NewRole(name string) *Role {
	return &Role{
		Name: name,
	}
}

// system roles should be saved in the db using initializer
type RoleName string

const (
	RoleAdmin RoleName = "Admin"
	RoleUser  RoleName = "User"
)