package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
)

func FromDomainPaddleCustomer(customer *domain.PaddleCustomer) *PaddleCustomerModel {
	if customer == nil {
		return nil
	}
	return &PaddleCustomerModel{ID: customer.ID, UserID: customer.UserID, PaddleCustomerID: customer.PaddleCustomerID}
}
