package domain

import "time"

type PlanItem struct {
	ID           string     `json:"id"`
	ImgUrl       *string    `json:"img_url"`
	ImgPath      *string    `json:"-"`
	Features     []string   `json:"features"`
	Status       PlanStatus `json:"status"`
	VarOptionIDs []string   `json:"var_option_ids,omitempty"`
	Options      []*Option  `json:"options"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type Option struct {
	Name        string `json:"name"`
	VarOptName  string `json:"var_opt_name"`
	VarOptValue string `json:"var_opt_value"`
}
