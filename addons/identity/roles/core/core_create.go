package core

type CreateRoleReq struct {
	Name string  `json:"name"`
}

func (r CreateRoleReq) Validate() error {
	return validatePostFields(r.Name)
}