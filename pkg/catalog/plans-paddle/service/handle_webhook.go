package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
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

func (s *Service) processPaddleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error {
	eventAt, err := parseWebhookTime(event.OccurredAt)
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

func (s *Service) handleSubscriptionWebhook(ctx context.Context, event *domain.PaddleWebhookEvent, eventAt time.Time) error {
	if event.Data.ID == "" {
		return fmt.Errorf("Paddle subscription webhook has no subscription ID")
	}
	if s.PaddleSubscriptionRepo == nil {
		return fmt.Errorf("Paddle subscription repository is not configured")
	}

	subscription := &domain.PaddleSubscription{
		PaddleSubscriptionID: event.Data.ID,
		UserID:               customDataString(event.Data.CustomData, "user_id"),
		PlanID:               customDataString(event.Data.CustomData, "plan_id"),
		PlanItemID:           customDataString(event.Data.CustomData, "plan_item_id"),
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
		subscription.CurrentPeriodStartsAt = optionalString(event.Data.CurrentBillingPeriod.StartsAt)
		subscription.CurrentPeriodEndsAt = optionalString(event.Data.CurrentBillingPeriod.EndsAt)
	}
	if event.Data.ScheduledChange != nil {
		subscription.ScheduledChange = &domain.PaddleSubscriptionChange{
			Action:      event.Data.ScheduledChange.Action,
			EffectiveAt: event.Data.ScheduledChange.EffectiveAt,
		}
	}

	if current, err := s.PaddleSubscriptionRepo.FindByPaddleSubscriptionID(ctx, subscription.PaddleSubscriptionID); err == nil {
		preserveSubscriptionFields(subscription, current)
	} else if !errors.Is(err, sharedD.ErrNotFound) {
		return fmt.Errorf("error loading Paddle subscription before webhook update: %w", err)
	}
	return s.PaddleSubscriptionRepo.Upsert(ctx, subscription, eventAt)
}

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
		UserID:              customDataString(event.Data.CustomData, "user_id"),
		PlanID:              customDataString(event.Data.CustomData, "plan_id"),
		PlanItemID:          customDataString(event.Data.CustomData, "plan_item_id"),
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
		transaction.Subtotal, err = parsePaddleAmount(event.Data.Details.Totals.Subtotal)
		if err != nil {
			return fmt.Errorf("invalid Paddle transaction subtotal: %w", err)
		}
		transaction.Tax, err = parsePaddleAmount(event.Data.Details.Totals.Tax)
		if err != nil {
			return fmt.Errorf("invalid Paddle transaction tax: %w", err)
		}
		transaction.Total, err = parsePaddleAmount(event.Data.Details.Totals.Total)
		if err != nil {
			return fmt.Errorf("invalid Paddle transaction total: %w", err)
		}
	}
	if existing, err := s.PaddleTransactionRepo.FindByPaddleTransactionID(ctx, transaction.PaddleTransactionID); err == nil {
		preserveTransactionFields(transaction, existing)
	} else if !errors.Is(err, sharedD.ErrNotFound) {
		return fmt.Errorf("error loading Paddle transaction before webhook update: %w", err)
	}
	return s.PaddleTransactionRepo.UpsertFromWebhook(ctx, transaction, eventAt)
}

func preserveSubscriptionFields(incoming, current *domain.PaddleSubscription) {
	if incoming.ID == "" {
		incoming.ID = current.ID
	}
	if incoming.PaddleSubscriptionID == "" {
		incoming.PaddleSubscriptionID = current.PaddleSubscriptionID
	}
	if incoming.UserID == "" {
		incoming.UserID = current.UserID
	}
	if incoming.PlanID == "" {
		incoming.PlanID = current.PlanID
	}
	if incoming.PlanItemID == "" {
		incoming.PlanItemID = current.PlanItemID
	}
	if incoming.CustomerID == "" {
		incoming.CustomerID = current.CustomerID
	}
}

func preserveTransactionFields(incoming, current *domain.PaddleTransaction) {
	if incoming.ID == "" {
		incoming.ID = current.ID
	}
	if incoming.PaddleTransactionID == "" {
		incoming.PaddleTransactionID = current.PaddleTransactionID
	}
	if incoming.UserID == "" {
		incoming.UserID = current.UserID
	}
	if incoming.CustomerID == "" {
		incoming.CustomerID = current.CustomerID
	}
	if incoming.PlanID == "" {
		incoming.PlanID = current.PlanID
	}
	if incoming.PlanItemID == "" {
		incoming.PlanItemID = current.PlanItemID
	}
	if incoming.SubscriptionID == nil {
		incoming.SubscriptionID = current.SubscriptionID
	}
	if current.CurrencyCode != "" {
		incoming.CurrencyCode = current.CurrencyCode
	}
	if current.Subtotal != 0 {
		incoming.Subtotal = current.Subtotal
	}
	if current.Tax != 0 {
		incoming.Tax = current.Tax
	}
	if current.Total != 0 {
		incoming.Total = current.Total
	}
	if current.PaymentDetailsProvided {
		incoming.PaymentDetailsProvided = true
	}
}

func (s *Service) resolveBillingCustomer(ctx context.Context, userID, paddleCustomerID string) (string, string, error) {
	if s.PaddleCustomerRepo == nil {
		return userID, "", nil
	}
	if paddleCustomerID != "" && userID == "" {
		customer, err := s.PaddleCustomerRepo.FindByPaddleCustomerID(ctx, paddleCustomerID)
		if err == nil {
			userID = customer.UserID
		} else if !errors.Is(err, sharedD.ErrNotFound) {
			return userID, "", fmt.Errorf("error finding GoEP user for Paddle customer: %w", err)
		}
	}
	if userID == "" {
		return userID, "", nil
	}
	if paddleCustomerID != "" {
		customer := &domain.PaddleCustomer{UserID: userID, PaddleCustomerID: paddleCustomerID}
		if err := s.PaddleCustomerRepo.Upsert(ctx, customer); err != nil {
			return userID, "", fmt.Errorf("error saving Paddle billing identity: %w", err)
		}
		return userID, customer.ID, nil
	}
	if customer, err := s.PaddleCustomerRepo.FindByUserID(ctx, userID); err == nil {
		return userID, customer.ID, nil
	} else if !errors.Is(err, sharedD.ErrNotFound) {
		return userID, "", fmt.Errorf("error finding Paddle billing identity: %w", err)
	}
	return userID, "", nil
}

func parsePaddleAmount(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	amount, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func customDataString(data map[string]any, key string) string {
	if data == nil {
		return ""
	}
	value, ok := data[key]
	if !ok {
		return ""
	}
	stringValue, _ := value.(string)
	return stringValue
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func parseWebhookTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339Nano, value)
}
