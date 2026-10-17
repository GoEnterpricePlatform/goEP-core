package model

import "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"

// PaddleCustomerModel stores the GoEP identity and Paddle billing ID separately.
type PaddleCustomerModel struct {
	ID               string `bson:"_id"`
	UserID           string `bson:"user_id"`
	PaddleCustomerID string `bson:"paddle_customer_id"`
}

func (m *PaddleCustomerModel) ToDomain() *domain.PaddleCustomer {
	if m == nil {
		return nil
	}
	return &domain.PaddleCustomer{ID: m.ID, PaddleCustomerID: m.PaddleCustomerID, UserID: m.UserID}
}

func FromDomainPaddleCustomer(customer *domain.PaddleCustomer) *PaddleCustomerModel {
	if customer == nil {
		return nil
	}
	return &PaddleCustomerModel{ID: customer.ID, UserID: customer.UserID, PaddleCustomerID: customer.PaddleCustomerID}
}
