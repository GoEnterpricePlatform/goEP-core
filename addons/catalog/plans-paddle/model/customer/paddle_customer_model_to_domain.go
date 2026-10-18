package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
)

func (m *PaddleCustomerModel) ToDomain() *domain.PaddleCustomer {
	if m == nil {
		return nil
	}
	return &domain.PaddleCustomer{ID: m.ID, PaddleCustomerID: m.PaddleCustomerID, UserID: m.UserID}
}
