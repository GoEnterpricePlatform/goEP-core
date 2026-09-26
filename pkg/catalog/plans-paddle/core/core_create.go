package core

type CreatePaddlePlanReq struct {
	Name            string                  `json:"name"`
	Description     *string                 `json:"description"`
	PaddleProductID string                  `json:"paddle_product_id"`
	Items           []*CreatePaddlePlanItem `json:"items"`
}

type CreatePaddlePlanItem struct {
	Features      []string `json:"features"`
	VarOptionIDs  []string `json:"var_option_ids,omitempty"`
	PaddlePriceID string   `json:"paddle_price_id"`
}
