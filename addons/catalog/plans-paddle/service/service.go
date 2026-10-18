package service

import (
	paddlePlanP "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/port"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

type Service struct {
	PaddlePlanTx           paddlePlanP.PaddlePlanTx
	PaddlePlanRepo         paddlePlanP.PaddlePlanRepo
	PaddleTransactionRepo  paddlePlanP.PaddleTransactionRepo
	PaddleSubscriptionRepo paddlePlanP.PaddleSubscriptionRepo
	PaddleCustomerRepo     paddlePlanP.PaddleCustomerRepo
	PaddleWebhookEventRepo paddlePlanP.PaddleWebhookEventRepo
	PlanRepo               port.PlanRepo
	PaddlePlanStg          paddlePlanP.PaddlePlanFileStg
	PaddleClient           *paddleSDK.SDK
}

func NewPlanPaddleSrv(paddlePlanTx paddlePlanP.PaddlePlanTx, paddlePlanRepo paddlePlanP.PaddlePlanRepo, paddleTransactionRepo paddlePlanP.PaddleTransactionRepo, paddleSubscriptionRepo paddlePlanP.PaddleSubscriptionRepo, planRepo port.PlanRepo, paddlePlanStg paddlePlanP.PaddlePlanFileStg, paddleClient *paddleSDK.SDK, paddleCustomerRepo paddlePlanP.PaddleCustomerRepo, paddleWebhookEventRepo paddlePlanP.PaddleWebhookEventRepo) *Service {
	return &Service{
		PaddlePlanTx:           paddlePlanTx,
		PaddlePlanRepo:         paddlePlanRepo,
		PaddleTransactionRepo:  paddleTransactionRepo,
		PaddleSubscriptionRepo: paddleSubscriptionRepo,
		PlanRepo:               planRepo,
		PaddlePlanStg:          paddlePlanStg,
		PaddleClient:           paddleClient,
		PaddleCustomerRepo:     paddleCustomerRepo,
		PaddleWebhookEventRepo: paddleWebhookEventRepo,
	}
}
