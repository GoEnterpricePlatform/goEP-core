package service

import (
	paddlePlanP "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/port"
)

type Service struct {
	PaddlePlanTx   paddlePlanP.PaddlePlanTx
	PaddlePlanRepo paddlePlanP.PaddlePlanRepo
	PlanRepo       port.PlanRepo
	PaddlePlanStg  paddlePlanP.PaddlePlanFileStg
}

func NewPlanPaddleSrv(paddlePlanTx paddlePlanP.PaddlePlanTx, paddlePlanRepo paddlePlanP.PaddlePlanRepo, planRepo port.PlanRepo, paddlePlanStg paddlePlanP.PaddlePlanFileStg) *Service {
	return &Service{
		PaddlePlanTx:   paddlePlanTx,
		PaddlePlanRepo: paddlePlanRepo,
		PlanRepo:       planRepo,
		PaddlePlanStg:  paddlePlanStg,
	}
}
