package service

import (
	"context"
	"testing"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/service/helpers"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

type memoryTransactionRepo struct {
	items map[string]*domain.PaddleTransaction
}

func (r *memoryTransactionRepo) Insert(_ context.Context, transaction *domain.PaddleTransaction) error {
	if r.items == nil {
		r.items = make(map[string]*domain.PaddleTransaction)
	}
	if transaction.ID == "" {
		transaction.ID = "internal:" + transaction.PaddleTransactionID
	}
	r.items[transaction.PaddleTransactionID] = transaction
	return nil
}

func (r *memoryTransactionRepo) FindByPaddleTransactionID(_ context.Context, id string) (*domain.PaddleTransaction, error) {
	if transaction := r.items[id]; transaction != nil {
		copy := *transaction
		return &copy, nil
	}
	return nil, sharedD.ErrNotFound
}

func (r *memoryTransactionRepo) FindAllByUserID(_ context.Context, userID string) ([]*domain.PaddleTransaction, error) {
	var result []*domain.PaddleTransaction
	for _, transaction := range r.items {
		if transaction.UserID == userID {
			copy := *transaction
			result = append(result, &copy)
		}
	}
	return result, nil
}

func (r *memoryTransactionRepo) FindAllLinkedToSubscriptions(_ context.Context) ([]*domain.PaddleTransaction, error) {
	var result []*domain.PaddleTransaction
	for _, transaction := range r.items {
		if transaction.SubscriptionID != nil {
			copy := *transaction
			result = append(result, &copy)
		}
	}
	return result, nil
}

func (r *memoryTransactionRepo) UpsertFromWebhook(_ context.Context, transaction *domain.PaddleTransaction, eventAt time.Time) error {
	if r.items == nil {
		r.items = make(map[string]*domain.PaddleTransaction)
	}
	if existing := r.items[transaction.PaddleTransactionID]; existing != nil && existing.LastEventAt.After(eventAt) {
		return nil
	}
	if existing := r.items[transaction.PaddleTransactionID]; existing != nil && transaction.ID == "" {
		transaction.ID = existing.ID
	}
	if transaction.ID == "" {
		transaction.ID = "internal:" + transaction.PaddleTransactionID
	}
	copy := *transaction
	copy.LastEventAt = eventAt
	r.items[transaction.PaddleTransactionID] = &copy
	return nil
}

type memorySubscriptionRepo struct {
	items map[string]*domain.PaddleSubscription
}

func (r *memorySubscriptionRepo) Upsert(_ context.Context, subscription *domain.PaddleSubscription, eventAt time.Time) error {
	if r.items == nil {
		r.items = make(map[string]*domain.PaddleSubscription)
	}
	if subscription.ID == "" {
		subscription.ID = "internal:" + subscription.PaddleSubscriptionID
	}
	if existing, ok := r.findByPaddleID(subscription.PaddleSubscriptionID); ok && existing.UpdatedAt != "" {
		if lastEventAt, err := helpers.ParseWebhookTime(existing.UpdatedAt); err == nil && lastEventAt.After(eventAt) {
			return nil
		}
	}
	copy := *subscription
	r.items[subscription.ID] = &copy
	return nil
}

func (r *memorySubscriptionRepo) findByPaddleID(id string) (*domain.PaddleSubscription, bool) {
	for _, subscription := range r.items {
		if subscription.PaddleSubscriptionID == id {
			return subscription, true
		}
	}
	return nil, false
}

func (r *memorySubscriptionRepo) FindByID(_ context.Context, id string) (*domain.PaddleSubscription, error) {
	if subscription := r.items[id]; subscription != nil {
		copy := *subscription
		return &copy, nil
	}
	return nil, sharedD.ErrNotFound
}

func (r *memorySubscriptionRepo) FindByPaddleSubscriptionID(_ context.Context, id string) (*domain.PaddleSubscription, error) {
	if subscription, ok := r.findByPaddleID(id); ok {
		copy := *subscription
		return &copy, nil
	}
	return nil, sharedD.ErrNotFound
}

func (r *memorySubscriptionRepo) FindCurrentByUserID(_ context.Context, userID string) (*domain.PaddleSubscription, error) {
	for _, subscription := range r.items {
		if subscription.UserID == userID && (subscription.Status == "active" || subscription.Status == "trialing") {
			copy := *subscription
			return &copy, nil
		}
	}
	return nil, sharedD.ErrNotFound
}

func (r *memorySubscriptionRepo) FindAllByUserID(_ context.Context, userID string) ([]*domain.PaddleSubscription, error) {
	var result []*domain.PaddleSubscription
	for _, subscription := range r.items {
		if subscription.UserID == userID {
			copy := *subscription
			result = append(result, &copy)
		}
	}
	return result, nil
}

func (r *memorySubscriptionRepo) FindAll(_ context.Context) ([]*domain.PaddleSubscription, error) {
	var result []*domain.PaddleSubscription
	for _, subscription := range r.items {
		copy := *subscription
		result = append(result, &copy)
	}
	return result, nil
}

func TestHandlePaddleWebhookSeparatesRenewalTransactionsFromSubscription(t *testing.T) {
	checkouts := &memoryTransactionRepo{items: map[string]*domain.PaddleTransaction{
		"txn_initial": {PaddleTransactionID: "txn_initial", UserID: "user-1", PlanID: "plan-1", PlanItemID: "item-1", Status: "draft"},
	}}
	subscriptions := &memorySubscriptionRepo{items: make(map[string]*domain.PaddleSubscription)}
	s := &Service{PaddleTransactionRepo: checkouts, PaddleSubscriptionRepo: subscriptions}

	created := &domain.PaddleWebhookEvent{
		EventType:  "subscription.created",
		OccurredAt: "2026-10-09T16:00:00Z",
		Data: domain.PaddleWebhookEventData{
			ID:     "sub_1",
			Status: "active",
			CustomData: map[string]any{
				"user_id":      "user-1",
				"plan_id":      "plan-1",
				"plan_item_id": "item-1",
			},
			CurrentBillingPeriod: &domain.PaddleWebhookTimePeriod{
				StartsAt: "2026-10-09T16:00:00Z",
				EndsAt:   "2026-11-09T16:00:00Z",
			},
		},
	}
	if err := s.HandlePaddleWebhook(context.Background(), created); err != nil {
		t.Fatalf("create subscription from webhook: %v", err)
	}

	for _, transactionID := range []string{"txn_initial", "txn_renewal"} {
		subscriptionID := "sub_1"
		event := &domain.PaddleWebhookEvent{
			EventType:  "transaction.completed",
			OccurredAt: "2026-10-09T16:01:00Z",
			Data: domain.PaddleWebhookEventData{
				ID:             transactionID,
				Status:         "completed",
				SubscriptionID: &subscriptionID,
				CustomData: map[string]any{
					"user_id":      "user-1",
					"plan_id":      "plan-1",
					"plan_item_id": "item-1",
				},
			},
		}
		if err := s.HandlePaddleWebhook(context.Background(), event); err != nil {
			t.Fatalf("save transaction %s: %v", transactionID, err)
		}
	}

	if got := len(subscriptions.items); got != 1 {
		t.Fatalf("subscription documents = %d, want 1", got)
	}
	if got := len(checkouts.items); got != 2 {
		t.Fatalf("transaction documents = %d, want 2", got)
	}
	for _, transactionID := range []string{"txn_initial", "txn_renewal"} {
		transaction := checkouts.items[transactionID]
		if transaction.SubscriptionID == nil || *transaction.SubscriptionID != "internal:sub_1" {
			t.Errorf("%s subscription_id = %v, want internal:sub_1", transactionID, transaction.SubscriptionID)
		}
		if transaction.Status != "completed" {
			t.Errorf("%s status = %q, want completed", transactionID, transaction.Status)
		}
	}
	if got := subscriptions.items["internal:sub_1"].Status; got != "active" {
		t.Errorf("subscription status = %q, want active", got)
	}
}

func TestGetSubscriptionAccessOnlyAllowsActiveAndTrialing(t *testing.T) {
	repo := &memorySubscriptionRepo{items: map[string]*domain.PaddleSubscription{
		"sub-active":   {ID: "sub-active", PaddleSubscriptionID: "paddle-active", UserID: "user-active", Status: "active"},
		"sub-trialing": {ID: "sub-trialing", PaddleSubscriptionID: "paddle-trialing", UserID: "user-trialing", Status: "trialing"},
		"sub-paused":   {ID: "sub-paused", PaddleSubscriptionID: "paddle-paused", UserID: "user-paused", Status: "paused"},
	}}
	s := &Service{PaddleSubscriptionRepo: repo}
	for userID, want := range map[string]bool{"user-active": true, "user-trialing": true, "user-paused": false, "user-none": false} {
		got, err := s.GetSubscriptionAccess(context.Background(), userID)
		if err != nil {
			t.Fatalf("access for %s: %v", userID, err)
		}
		if got.HasAccess != want {
			t.Errorf("access for %s = %v, want %v", userID, got.HasAccess, want)
		}
	}
}
