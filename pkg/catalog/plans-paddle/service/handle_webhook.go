package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
)

func (s *Service) HandlePaddleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error {
	if event == nil {
		return fmt.Errorf("Paddle webhook event is nil")
	}
	if event.Data.ID == "" {
		return fmt.Errorf("Paddle webhook has no transaction ID")
	}

	status := event.Data.Status
	switch event.EventType {
	case "transaction.completed":
		status = "completed"
	case "transaction.canceled", "transaction.past_due":
		if status == "" {
			status = event.EventType[len("transaction."):]
		}
	default:
		return nil
	}
	if status == "" {
		return fmt.Errorf("Paddle webhook has no transaction status")
	}
	return s.PaddleCheckoutRepo.UpdateStatus(ctx, event.Data.ID, status, event.Data.SubscriptionID, time.Now().UTC().Format(time.RFC3339Nano))
}
