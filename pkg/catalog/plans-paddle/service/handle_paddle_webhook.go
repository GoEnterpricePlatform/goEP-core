package service

import (
	"context"
	"fmt"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"time"
)

func (s *Service) HandlePaddleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error {
	if event == nil {
		return fmt.Errorf("Paddle webhook event is nil")
	}
	if event.EventID == "" {
		return fmt.Errorf("Paddle webhook event has no event ID")
	}
	if s.PaddleWebhookEventRepo != nil {
		alreadyProcessed, err := s.PaddleWebhookEventRepo.StartProcessing(ctx, event, time.Now().UTC())
		if err != nil {
			return fmt.Errorf("error starting Paddle webhook event %s: %w", event.EventID, err)
		}
		if alreadyProcessed {
			return nil
		}
	}

	if err := s.processPaddleWebhook(ctx, event); err != nil {
		if s.PaddleWebhookEventRepo != nil {
			if auditErr := s.PaddleWebhookEventRepo.MarkFailed(ctx, event.EventID, time.Now().UTC(), err); auditErr != nil {
				return fmt.Errorf("%w; error recording failed Paddle webhook event: %v", err, auditErr)
			}
		}
		return err
	}
	if s.PaddleWebhookEventRepo != nil {
		if err := s.PaddleWebhookEventRepo.MarkProcessed(ctx, event.EventID, time.Now().UTC()); err != nil {
			return fmt.Errorf("error marking Paddle webhook event %s as processed: %w", event.EventID, err)
		}
	}
	return nil
}
