package service

import (
	billingPort "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/port"
	catalogPaddlePort "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

type Service struct {
	PaddleTransactionRepo  billingPort.PaddleTransactionRepo
	PaddleSubscriptionRepo billingPort.PaddleSubscriptionRepo
	PaddleCustomerRepo     billingPort.PaddleCustomerRepo
	PaddleWebhookEventRepo billingPort.PaddleWebhookEventRepo
	CatalogCheckoutSrv     catalogPaddlePort.PaddleCheckoutCatalogSrv
	PaddleClient           *paddleSDK.SDK
}

var _ billingPort.PaddleBillingSrv = (*Service)(nil)

func NewPaddleBillingSrv(
	transactions billingPort.PaddleTransactionRepo,
	subscriptions billingPort.PaddleSubscriptionRepo,
	customers billingPort.PaddleCustomerRepo,
	webhookEvents billingPort.PaddleWebhookEventRepo,
	catalogCheckoutSrv catalogPaddlePort.PaddleCheckoutCatalogSrv,
	paddleClient *paddleSDK.SDK,
) *Service {
	return &Service{
		PaddleTransactionRepo:  transactions,
		PaddleSubscriptionRepo: subscriptions,
		PaddleCustomerRepo:     customers,
		PaddleWebhookEventRepo: webhookEvents,
		CatalogCheckoutSrv:     catalogCheckoutSrv,
		PaddleClient:           paddleClient,
	}
}
