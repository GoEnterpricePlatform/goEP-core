package domain

// PaddleCustomer links a GoEP user to the corresponding Paddle billing identity.
type PaddleCustomer struct {
	ID               string `json:"id"`
	PaddleCustomerID string `json:"paddle_customer_id"`
	UserID           string `json:"user_id"`
}
