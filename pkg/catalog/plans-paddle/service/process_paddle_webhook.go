package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/service/helpers"
)

func (s *Service) processPaddleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error {
	eventAt, err := helpers.ParseWebhookTime(event.OccurredAt)
	if err != nil {
		return fmt.Errorf("Paddle webhook has invalid occurred_at: %w", err)
	}
	if eventAt.IsZero() {
		eventAt = time.Now().UTC()
	}

	switch event.EventType {
	case "subscription.created", "subscription.updated", "subscription.activated", "subscription.canceled", "subscription.paused", "subscription.resumed", "subscription.past_due", "subscription.trialing":
		return s.handleSubscriptionWebhook(ctx, event, eventAt)
	case "transaction.completed", "transaction.canceled", "transaction.past_due", "transaction.payment_failed":
		return s.handleTransactionWebhook(ctx, event, eventAt)
	default:
		// Paddle sends unrelated events to the same notification destination.
		return nil
	}
}
