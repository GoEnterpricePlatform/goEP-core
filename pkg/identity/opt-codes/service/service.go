package service

import port "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/port"

var _ port.OtpCodeSrv = &Service{}

type Service struct {
	OtpCodeRepo port.OtpCodeRepo
}

func NewOtpCodeSrv(otpCodeRepo port.OtpCodeRepo) *Service {
	return &Service{
		OtpCodeRepo: otpCodeRepo,
	}
}
