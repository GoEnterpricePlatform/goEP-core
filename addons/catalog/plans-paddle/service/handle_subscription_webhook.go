package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/service/helpers"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

func (s *Service) handleSubscriptionWebhook(ctx context.Context, event *domain.PaddleWebhookEvent, eventAt time.Time) error {
	if event.Data.ID == "" {
		return fmt.Errorf("Paddle subscription webhook has no subscription ID")
	}
	if s.PaddleSubscriptionRepo == nil {
		return fmt.Errorf("Paddle subscription repository is not configured")
	}

	subscription := &domain.PaddleSubscription{
		PaddleSubscriptionID: event.Data.ID,
		UserID:               helpers.CustomDataString(event.Data.CustomData, "user_id"),
		PlanID:               helpers.CustomDataString(event.Data.CustomData, "plan_id"),
		PlanItemID:           helpers.CustomDataString(event.Data.CustomData, "plan_item_id"),
		Status:               event.Data.Status,
		NextBilledAt:         event.Data.NextBilledAt,
		CanceledAt:           event.Data.CanceledAt,
		UpdatedAt:            event.Data.UpdatedAt,
	}
	resolvedUserID, customerID, err := s.resolveBillingCustomer(ctx, subscription.UserID, event.Data.CustomerID)
	if err != nil {
		return err
	}
	subscription.UserID = resolvedUserID
	subscription.CustomerID = customerID
	if subscription.UpdatedAt == "" {
		subscription.UpdatedAt = event.OccurredAt
	}
	if subscription.Status == "" {
		return fmt.Errorf("Paddle subscription webhook has no status")
	}
	if event.Data.CurrentBillingPeriod != nil {
		subscription.CurrentPeriodStartsAt = helpers.OptionalString(event.Data.CurrentBillingPeriod.StartsAt)
		subscription.CurrentPeriodEndsAt = helpers.OptionalString(event.Data.CurrentBillingPeriod.EndsAt)
	}
	if event.Data.ScheduledChange != nil {
		subscription.ScheduledChange = &domain.PaddleSubscriptionChange{
			Action:      event.Data.ScheduledChange.Action,
			EffectiveAt: event.Data.ScheduledChange.EffectiveAt,
		}
	}

	if current, err := s.PaddleSubscriptionRepo.FindByPaddleSubscriptionID(ctx, subscription.PaddleSubscriptionID); err == nil {
		helpers.PreserveSubscriptionFields(subscription, current)
	} else if !errors.Is(err, sharedD.ErrNotFound) {
		return fmt.Errorf("error loading Paddle subscription before webhook update: %w", err)
	}
	return s.PaddleSubscriptionRepo.Upsert(ctx, subscription, eventAt)
}
