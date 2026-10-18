package model

// PaddleCustomerModel stores the GoEP identity and Paddle billing ID separately.
type PaddleCustomerModel struct {
	ID               string `bson:"_id"`
	UserID           string `bson:"user_id"`
	PaddleCustomerID string `bson:"paddle_customer_id"`
}
