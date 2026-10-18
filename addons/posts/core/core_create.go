package core

type CreatePostReq struct {
	Title string  `json:"title"`
	Desc *string `json:"desc"`
}

func (p CreatePostReq) Validate() error {
	return validatePostFields(p.Title, p.Desc)
}