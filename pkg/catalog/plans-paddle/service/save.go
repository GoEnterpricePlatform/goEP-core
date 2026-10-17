package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

// ReconcileUserSubscriptions refreshes the authenticated user's cached Paddle
// subscription state from Paddle's API. Normal entitlement checks remain local.
func (s *Service) ReconcileUserSubscriptions(ctx context.Context, userID string) error {
	if userID == "" {
		return sharedD.NewAppError(sharedD.ErrCodeUnauthorized, "authentication is required")
	}
	if s.PaddleClient == nil {
		return sharedD.NewAppError(sharedD.ErrCodeInternalServerError, "Paddle is not configured")
	}
	subscriptions, err := s.PaddleSubscriptionRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return sharedD.ManageError(err, "error finding Paddle subscriptions for reconciliation")
	}
	localByID := make(map[string]*domain.PaddleSubscription, len(subscriptions))
	for _, subscription := range subscriptions {
		localByID[subscription.ID] = subscription
	}
	checkouts, err := s.PaddleTransactionRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return sharedD.ManageError(err, "error finding Paddle checkouts for reconciliation")
	}
	for _, checkout := range checkouts {
		if checkout.SubscriptionID != nil && *checkout.SubscriptionID != "" {
			continue
		}
		remoteTransaction, err := s.PaddleClient.GetTransaction(ctx, &paddleSDK.GetTransactionRequest{TransactionID: checkout.PaddleTransactionID})
		if err != nil {
			return fmt.Errorf("error reconciling Paddle transaction %s: %w", checkout.PaddleTransactionID, err)
		}
		if remoteTransaction == nil || remoteTransaction.SubscriptionID == nil || *remoteTransaction.SubscriptionID == "" {
			continue
		}
		if subscription, err := s.PaddleSubscriptionRepo.FindByPaddleSubscriptionID(ctx, *remoteTransaction.SubscriptionID); err == nil {
			checkout.SubscriptionID = &subscription.ID
		} else if !errors.Is(err, sharedD.ErrNotFound) {
			return fmt.Errorf("error finding linked Paddle subscription: %w", err)
		}
		checkout.Status = string(remoteTransaction.Status)
		checkout.UpdatedAt = remoteTransaction.UpdatedAt
		if checkout.UserID == "" {
			checkout.UserID = customDataString(remoteTransaction.CustomData, "user_id")
		}
		if checkout.PlanID == "" {
			checkout.PlanID = customDataString(remoteTransaction.CustomData, "plan_id")
		}
		if checkout.PlanItemID == "" {
			checkout.PlanItemID = customDataString(remoteTransaction.CustomData, "plan_item_id")
		}
		transactionUpdatedAt, parseErr := parseWebhookTime(remoteTransaction.UpdatedAt)
		if parseErr != nil || transactionUpdatedAt.IsZero() {
			transactionUpdatedAt = time.Now().UTC()
			checkout.UpdatedAt = transactionUpdatedAt.Format(time.RFC3339Nano)
		}
		if err := s.PaddleTransactionRepo.UpsertFromWebhook(ctx, checkout, transactionUpdatedAt); err != nil {
			return fmt.Errorf("error saving reconciled Paddle transaction %s: %w", checkout.PaddleTransactionID, err)
		}
	}
	checkouts, err = s.PaddleTransactionRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return sharedD.ManageError(err, "error reloading Paddle checkouts for reconciliation")
	}
	for _, checkout := range checkouts {
		if checkout.SubscriptionID != nil && *checkout.SubscriptionID != "" {
			if _, exists := localByID[*checkout.SubscriptionID]; !exists {
				if subscription, err := s.PaddleSubscriptionRepo.FindByID(ctx, *checkout.SubscriptionID); err == nil {
					localByID[subscription.ID] = subscription
				} else if !errors.Is(err, sharedD.ErrNotFound) {
					return fmt.Errorf("error finding linked Paddle subscription: %w", err)
				}
			}
		}
	}
	for _, local := range localByID {
		remote, err := s.PaddleClient.GetSubscription(ctx, &paddleSDK.GetSubscriptionRequest{SubscriptionID: local.PaddleSubscriptionID})
		if err != nil {
			return fmt.Errorf("error reconciling Paddle subscription %s: %w", local.PaddleSubscriptionID, err)
		}
		if remote == nil || remote.ID == "" {
			return fmt.Errorf("Paddle returned an empty subscription for %s", local.PaddleSubscriptionID)
		}
		updatedAt := remote.UpdatedAt
		if parsedUpdatedAt, parseErr := parseWebhookTime(updatedAt); parseErr != nil || parsedUpdatedAt.IsZero() {
			updatedAt = time.Now().UTC().Format(time.RFC3339Nano)
		}
		// A reconciliation fetch is authoritative even if Paddle's object
		// updated_at is slightly earlier than a webhook's occurred_at.
		lastEventAt := time.Now().UTC()
		fresh := &domain.PaddleSubscription{
			ID:                   local.ID,
			PaddleSubscriptionID: remote.ID,
			UserID:               customDataString(remote.CustomData, "user_id"),
			CustomerID:           local.CustomerID,
			PlanID:               customDataString(remote.CustomData, "plan_id"),
			PlanItemID:           customDataString(remote.CustomData, "plan_item_id"),
			Status:               string(remote.Status),
			NextBilledAt:         remote.NextBilledAt,
			CanceledAt:           remote.CanceledAt,
			UpdatedAt:            updatedAt,
		}
		if fresh.UserID == "" {
			fresh.UserID = userID
		}
		if fresh.PlanID == "" {
			fresh.PlanID = local.PlanID
		}
		if fresh.PlanItemID == "" {
			fresh.PlanItemID = local.PlanItemID
		}
		if remote.CurrentBillingPeriod != nil {
			fresh.CurrentPeriodStartsAt = optionalString(remote.CurrentBillingPeriod.StartsAt)
			fresh.CurrentPeriodEndsAt = optionalString(remote.CurrentBillingPeriod.EndsAt)
		}
		if remote.ScheduledChange != nil {
			fresh.ScheduledChange = &domain.PaddleSubscriptionChange{
				Action:      string(remote.ScheduledChange.Action),
				EffectiveAt: remote.ScheduledChange.EffectiveAt,
			}
		}
		if err := s.PaddleSubscriptionRepo.Upsert(ctx, fresh, lastEventAt); err != nil {
			return fmt.Errorf("error saving reconciled Paddle subscription %s: %w", local.PaddleSubscriptionID, err)
		}
	}
	return nil
}
