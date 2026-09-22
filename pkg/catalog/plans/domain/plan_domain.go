package domain

import "time"

type Plan struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	ImgUrl      *string      `json:"img_url"`
	ImgPath     *string      `json:"-"`
	Items       []*PlanItem  `json:"items,omitempty"`
	Variations  []*Variation `json:"variations"`
	CreatedAt   *time.Time   `json:"created_at"`
	UpdatedAt   *time.Time   `json:"updated_at"`
}

// Only response, It is assigned from service with the calculate_variations function, therefore it does not need a structure in model
type Variation struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}


