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
