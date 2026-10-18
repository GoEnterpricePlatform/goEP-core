package service

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/billing/payment-providers/port"
	encryptor "github.com/GoEnterpricePlatform/goEP-core/internal/shared/encryptor"
)

var _ port.PaymentProviderSrv = &Service{}

type Service struct {
	PProviderRepo port.PaymentProviderRepo
	EncryptorSrv encryptor.EncryptorSrv
}

func NewPaymentProviderSrv(pProviderRepo port.PaymentProviderRepo, encryptorSrv encryptor.EncryptorSrv) *Service {
	return &Service{
		PProviderRepo: pProviderRepo,
		EncryptorSrv: encryptorSrv,
	}
}
