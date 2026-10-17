package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

func (s *Service) GetSubscriptionAccess(ctx context.Context, userID string) (*domain.PaddleSubscriptionAccess, error) {
	if userID == "" {
		return nil, sharedD.NewAppError(sharedD.ErrCodeUnauthorized, "authentication is required")
	}
	subscription, err := s.PaddleSubscriptionRepo.FindCurrentByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sharedD.ErrNotFound) {
			return &domain.PaddleSubscriptionAccess{HasAccess: false}, nil
		}
		return nil, sharedD.ManageError(err, "error finding current Paddle subscription")
	}
	return &domain.PaddleSubscriptionAccess{HasAccess: true, Subscription: subscription}, nil
}

func (s *Service) RunSubscriptionReconciliation(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 6 * time.Hour
	}
	if err := s.ReconcileAllSubscriptions(ctx); err != nil {
		log.Printf("initial Paddle subscription reconciliation failed: %v", err)
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.ReconcileAllSubscriptions(ctx); err != nil {
				log.Printf("Paddle subscription reconciliation failed: %v", err)
			}
		}
	}
}

func (s *Service) ReconcileAllSubscriptions(ctx context.Context) error {
	if s.PaddleClient == nil {
		return fmt.Errorf("Paddle is not configured")
	}
	subscriptions, err := s.PaddleSubscriptionRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error listing Paddle subscriptions for reconciliation: %w", err)
	}
	byID := make(map[string]*domain.PaddleSubscription, len(subscriptions))
	for _, subscription := range subscriptions {
		byID[subscription.ID] = subscription
	}
	for _, local := range byID {
		remote, err := s.PaddleClient.GetSubscription(ctx, &paddleSDK.GetSubscriptionRequest{SubscriptionID: local.PaddleSubscriptionID})
		if err != nil {
			return fmt.Errorf("error reconciling Paddle subscription %s: %w", local.PaddleSubscriptionID, err)
		}
		if remote == nil || remote.ID == "" {
			return fmt.Errorf("Paddle returned an empty subscription for %s", local.PaddleSubscriptionID)
		}
		updatedAt := remote.UpdatedAt
		if parsed, parseErr := parseWebhookTime(updatedAt); parseErr != nil || parsed.IsZero() {
			updatedAt = time.Now().UTC().Format(time.RFC3339Nano)
		}
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
			fresh.UserID = local.UserID
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
		if err := s.PaddleSubscriptionRepo.Upsert(ctx, fresh, time.Now().UTC()); err != nil {
			return fmt.Errorf("error saving reconciled Paddle subscription %s: %w", local.PaddleSubscriptionID, err)
		}
	}
	return nil
}
