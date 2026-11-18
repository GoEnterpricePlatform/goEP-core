package core

type UpdatePostReq struct {
	Title string  `json:"title"`
	Desc *string `json:"desc"`
}

func (u UpdatePostReq) Validate() error {
	return validatePostFields(u.Title, u.Desc)
}