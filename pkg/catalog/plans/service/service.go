package service

import (
	planP "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/port"
)

var _ planP.PlanSrv = &Service{}

type Service struct {
	PlanRepo      planP.PlanRepo
	PlanFileStg   planP.PlanFileStg
	VarOptionRepo port.VarOptionRepo
}

func NewPlanSrv(PlanRepo planP.PlanRepo, planFileStg planP.PlanFileStg, varOptionRepo port.VarOptionRepo) *Service {
	return &Service{
		PlanRepo:      PlanRepo,
		PlanFileStg:   planFileStg,
		VarOptionRepo: varOptionRepo,
	}
}
