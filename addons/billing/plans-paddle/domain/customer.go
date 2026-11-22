package domain

// PaddleCustomer is Paddle's billing identity linked to a GoEP user.
type PaddleCustomer struct {
	ID               string `json:"id"`
	PaddleCustomerID string `json:"paddle_customer_id"`
	UserID           string `json:"user_id"`
}
