package service

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
)

func (s *Service) HandlePaddleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error {
	if event == nil {
		return fmt.Errorf("Paddle webhook event is nil")
	}
	if event.EventID == "" {
		return fmt.Errorf("Paddle webhook event has no event ID")
	}
	switch event.EventType {
	case "product.created", "product.updated":
		return s.handleCatalogPaddleWebhook(ctx, event)
	default:
		return fmt.Errorf("unsupported Catalog Paddle event type %q", event.EventType)
	}
}
