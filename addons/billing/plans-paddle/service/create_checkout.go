package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

func (s *Service) CreateCheckout(ctx context.Context, userID, paddlePlanID, paddlePlanItemID string) (*domain.PaddleTransaction, error) {
	if userID == "" {
		return nil, sharedD.NewAppError(sharedD.ErrCodeUnauthorized, "authentication is required to create a Paddle checkout")
	}
	if paddlePlanID == "" || paddlePlanItemID == "" {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "Paddle plan and plan item are required")
	}
	if s.CatalogCheckoutSrv == nil {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInternalServerError, "Catalog checkout service is not configured")
	}
	offer, err := s.CatalogCheckoutSrv.ResolvePaddleCheckout(ctx, paddlePlanID, paddlePlanItemID)
	if err != nil {
		return nil, err
	}
	if s.PaddleClient == nil {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInternalServerError, "Paddle is not configured")
	}
	transactionItem := paddleSDK.NewCreateTransactionItemsTransactionItemFromCatalog(
		&paddleSDK.TransactionItemFromCatalog{PriceID: offer.PaddlePriceID, Quantity: 1},
	)
	transaction, err := s.PaddleClient.CreateTransaction(ctx, &paddleSDK.CreateTransactionRequest{
		Items: []paddleSDK.CreateTransactionItems{*transactionItem},
		CustomData: paddleSDK.CustomData{
			"plan_id":      offer.PlanID,
			"plan_item_id": offer.PlanItemID,
			"user_id":      userID,
		},
	})
	if err != nil {
		return nil, sharedD.ManageError(err, "error creating Paddle transaction")
	}
	if transaction == nil || transaction.ID == "" {
		return nil, fmt.Errorf("Paddle returned an empty transaction ID")
	}

	paddleTransaction := &domain.PaddleTransaction{
		PaddleTransactionID: transaction.ID,
		UserID:              userID,
		PlanID:              offer.PlanID,
		PlanItemID:          offer.PlanItemID,
		Status:              string(transaction.Status),
		UpdatedAt:           time.Now().UTC().Format(time.RFC3339Nano),
	}
	if s.PaddleCustomerRepo != nil {
		if customer, err := s.PaddleCustomerRepo.FindByUserID(ctx, userID); err == nil {
			paddleTransaction.CustomerID = customer.ID
		} else if !errors.Is(err, sharedD.ErrNotFound) {
			return nil, sharedD.ManageError(err, "error finding billing customer")
		}
	}
	if err := s.PaddleTransactionRepo.Insert(ctx, paddleTransaction); err != nil {
		return nil, sharedD.ManageError(err, "error saving Paddle checkout")
	}
	return paddleTransaction, nil
}
