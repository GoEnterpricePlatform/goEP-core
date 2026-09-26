package service

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
)

type Service struct {
	PaddlePlanTx port.PaddlePlanTx
}

func NewPlanPaddleSrv(paddlePlanTx port.PaddlePlanTx) *Service {
	return &Service{
		PaddlePlanTx: paddlePlanTx,
	}
}
