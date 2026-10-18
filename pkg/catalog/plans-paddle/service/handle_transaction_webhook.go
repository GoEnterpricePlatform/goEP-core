package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/service/helpers"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) handleTransactionWebhook(ctx context.Context, event *domain.PaddleWebhookEvent, eventAt time.Time) error {
	if event.Data.ID == "" {
		return fmt.Errorf("Paddle transaction webhook has no transaction ID")
	}
	status := event.Data.Status
	switch event.EventType {
	case "transaction.completed":
		status = "completed"
	case "transaction.canceled":
		status = "canceled"
	case "transaction.past_due", "transaction.payment_failed":
		status = "past_due"
	}
	if status == "" {
		return fmt.Errorf("Paddle transaction webhook has no status")
	}

	transaction := &domain.PaddleTransaction{
		PaddleTransactionID: event.Data.ID,
		UserID:              helpers.CustomDataString(event.Data.CustomData, "user_id"),
		PlanID:              helpers.CustomDataString(event.Data.CustomData, "plan_id"),
		PlanItemID:          helpers.CustomDataString(event.Data.CustomData, "plan_item_id"),
		Status:              status,
		UpdatedAt:           eventAt.Format(time.RFC3339Nano),
		LastEventAt:         eventAt,
	}
	resolvedUserID, customerID, err := s.resolveBillingCustomer(ctx, transaction.UserID, event.Data.CustomerID)
	if err != nil {
		return err
	}
	transaction.UserID = resolvedUserID
	transaction.CustomerID = customerID
	if event.Data.SubscriptionID != nil && *event.Data.SubscriptionID != "" && s.PaddleSubscriptionRepo != nil {
		if subscription, err := s.PaddleSubscriptionRepo.FindByPaddleSubscriptionID(ctx, *event.Data.SubscriptionID); err == nil {
			transaction.SubscriptionID = &subscription.ID
			if transaction.UserID == "" {
				transaction.UserID = subscription.UserID
			}
			if transaction.CustomerID == "" {
				transaction.CustomerID = subscription.CustomerID
			}
			if transaction.PlanID == "" {
				transaction.PlanID = subscription.PlanID
			}
			if transaction.PlanItemID == "" {
				transaction.PlanItemID = subscription.PlanItemID
			}
		} else if !errors.Is(err, sharedD.ErrNotFound) {
			return fmt.Errorf("error finding related Paddle subscription: %w", err)
		}
	}
	if event.Data.Details != nil {
		transaction.PaymentDetailsProvided = true
		transaction.CurrencyCode = event.Data.CurrencyCode
		transaction.Subtotal, err = helpers.ParsePaddleAmount(event.Data.Details.Totals.Subtotal)
		if err != nil {
			return fmt.Errorf("invalid Paddle transaction subtotal: %w", err)
		}
		transaction.Tax, err = helpers.ParsePaddleAmount(event.Data.Details.Totals.Tax)
		if err != nil {
			return fmt.Errorf("invalid Paddle transaction tax: %w", err)
		}
		transaction.Total, err = helpers.ParsePaddleAmount(event.Data.Details.Totals.Total)
		if err != nil {
			return fmt.Errorf("invalid Paddle transaction total: %w", err)
		}
	}
	if existing, err := s.PaddleTransactionRepo.FindByPaddleTransactionID(ctx, transaction.PaddleTransactionID); err == nil {
		helpers.PreserveTransactionFields(transaction, existing)
	} else if !errors.Is(err, sharedD.ErrNotFound) {
		return fmt.Errorf("error loading Paddle transaction before webhook update: %w", err)
	}
	return s.PaddleTransactionRepo.UpsertFromWebhook(ctx, transaction, eventAt)
}
