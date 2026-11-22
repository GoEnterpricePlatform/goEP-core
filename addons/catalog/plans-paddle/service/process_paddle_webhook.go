package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)

// handleCatalogPaddleWebhook mirrors Paddle product names and descriptions to
// the linked GoEP plan. Paddle product details remain authoritative for this
// Paddle-backed plan mapping.
func (s *Service) handleCatalogPaddleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error {
	product := event.Data
	if product.ID == "" {
		return fmt.Errorf("Paddle product webhook has no product ID")
	}
	if strings.TrimSpace(product.Name) == "" {
		return fmt.Errorf("Paddle product webhook has no product name")
	}

	paddlePlan, err := s.PaddlePlanRepo.FindByPaddleProductID(ctx, product.ID)
	if err != nil {
		if errors.Is(err, sharedD.ErrNotFound) {
			// Products created in Paddle can arrive before an administrator links
			// them to a GoEP plan. There is nothing local to update in that case.
			return nil
		}
		return fmt.Errorf("find Catalog mapping for Paddle product %s: %w", product.ID, err)
	}

	plan, err := s.PlanRepo.Find(ctx, paddlePlan.PlanID)
	if err != nil {
		return fmt.Errorf("find GoEP plan linked to Paddle product %s: %w", product.ID, err)
	}
	if plan.Name == product.Name && sameOptionalString(plan.Description, product.Description) {
		return nil
	}
	plan.Name = product.Name
	plan.Description = product.Description
	now := time.Now().UTC()
	plan.UpdatedAt = &now
	if err := s.PlanRepo.Update(ctx, plan.ID, plan); err != nil {
		return fmt.Errorf("sync GoEP plan from Paddle product %s: %w", product.ID, err)
	}
	return nil
}

func sameOptionalString(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
