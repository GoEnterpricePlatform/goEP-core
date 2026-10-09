package service

import (
	"context"
	"fmt"
	"time"

	paddlePlanD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

func (s *Service) CreateCheckout(ctx context.Context, planID, itemID string) (*paddlePlanD.PaddleCheckout, error) {
	if s.PaddleClient == nil {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInternalServerError, "Paddle is not configured")
	}
	plan, err := s.PaddlePlanRepo.Find(ctx, planID)
	if err != nil {
		return nil, sharedD.ManageError(err, "error finding plan for checkout")
	}

	var selected *paddlePlanD.PaddlePlanItem
	for _, item := range plan.Items {
		if item != nil && item.ID == itemID {
			selected = item

			break
		}
	}
	if selected == nil {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "plan item not found")
	}

	// The item's canonical status is stored in catalog_plans.items. Existing
	// catalog_plans_paddle documents only contain the Paddle price mapping.
	catalogPlan, err := s.PlanRepo.Find(ctx, plan.PlanID)
	if err != nil {
		return nil, sharedD.ManageError(err, "error finding catalog plan for checkout")
	}
	var itemIsActive bool
	for _, catalogItem := range catalogPlan.Items {
		if catalogItem != nil && catalogItem.ID == selected.PlanItemID {
			itemIsActive = string(catalogItem.Status) == "active"
			break
		}
	}

	if !itemIsActive {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "plan item is not available for purchase")
	}
	if selected.PaddlePriceID == "" {
		return nil, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "plan item has no Paddle price")
	}

	transactionItem := paddleSDK.NewCreateTransactionItemsTransactionItemFromCatalog(
		&paddleSDK.TransactionItemFromCatalog{PriceID: selected.PaddlePriceID, Quantity: 1},
	)
	transaction, err := s.PaddleClient.CreateTransaction(ctx, &paddleSDK.CreateTransactionRequest{
		Items: []paddleSDK.CreateTransactionItems{*transactionItem},
		CustomData: paddleSDK.CustomData{
			"plan_id":      plan.PlanID,
			"plan_item_id": selected.PlanItemID,
		},
	})
	if err != nil {
		return nil, sharedD.ManageError(err, "error creating Paddle transaction")
	}
	if transaction == nil || transaction.ID == "" {
		return nil, fmt.Errorf("Paddle returned an empty transaction ID")
	}

	checkout := &paddlePlanD.PaddleCheckout{
		TransactionID: transaction.ID,
		PlanID:        plan.PlanID,
		PlanItemID:    selected.PlanItemID,
		PaddlePriceID: selected.PaddlePriceID,
		Status:        string(transaction.Status),
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339Nano),
	}
	if err := s.PaddleCheckoutRepo.Insert(ctx, checkout); err != nil {
		return nil, sharedD.ManageError(err, "error saving Paddle checkout")
	}

	return checkout, nil
}
