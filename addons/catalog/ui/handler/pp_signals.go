package handler

type paddlePlanSignals struct {
	PlanForm struct {
		ID                string                  `json:"id"`
		Name              string                  `json:"name"`
		Description       string                  `json:"description"`
		PaddleProductID   string                  `json:"paddle_product_id"`
		Order             int                     `json:"order"`
		SelectedOptionIDs []string                `json:"selected_option_ids"`
		ProductType       string                  `json:"product_type"`
		Items             []paddlePlanItemSignals `json:"items"`
	} `json:"plan_form"`
}

type paddlePlanItemSignals struct {
	ID            string   `json:"id"`
	VarOptionIDs  []string `json:"var_option_ids"`
	Features      []string `json:"features"`
	PaddlePriceID string   `json:"paddle_price_id"`
	Status        string   `json:"status"`
}
