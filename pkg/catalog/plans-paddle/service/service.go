package service

import (
	paddlePlanP "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/port"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

type Service struct {
	PaddlePlanTx       paddlePlanP.PaddlePlanTx
	PaddlePlanRepo     paddlePlanP.PaddlePlanRepo
	PaddleCheckoutRepo paddlePlanP.PaddleCheckoutRepo
	PlanRepo           port.PlanRepo
	PaddlePlanStg      paddlePlanP.PaddlePlanFileStg
	PaddleClient       *paddleSDK.SDK
}

func NewPlanPaddleSrv(paddlePlanTx paddlePlanP.PaddlePlanTx, paddlePlanRepo paddlePlanP.PaddlePlanRepo, paddleCheckoutRepo paddlePlanP.PaddleCheckoutRepo, planRepo port.PlanRepo, paddlePlanStg paddlePlanP.PaddlePlanFileStg, paddleClient *paddleSDK.SDK) *Service {
	return &Service{
		PaddlePlanTx:       paddlePlanTx,
		PaddlePlanRepo:     paddlePlanRepo,
		PaddleCheckoutRepo: paddleCheckoutRepo,
		PlanRepo:           planRepo,
		PaddlePlanStg:      paddlePlanStg,
		PaddleClient:       paddleClient,
	}
}
