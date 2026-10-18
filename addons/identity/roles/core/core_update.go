package core

type UpdateRoleReq struct {
	Name string  `json:"name"`
}

func (u UpdateRoleReq) Validate() error {
	return validatePostFields(u.Name)
}