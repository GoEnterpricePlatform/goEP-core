package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/claim"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
)

type checkoutServiceStub struct {
	userID string
	called bool
}

func (s *checkoutServiceStub) Create(context.Context, *domain.PaddlePlan) error { return nil }
func (s *checkoutServiceStub) Get(context.Context, string) (*domain.PaddlePlan, error) {
	return nil, nil
}
func (s *checkoutServiceStub) GetAll(context.Context, int64, int64) ([]*domain.PaddlePlan, int64, int64, error) {
	return nil, 0, 0, nil
}
func (s *checkoutServiceStub) Update(context.Context, string, *domain.PaddlePlan) error { return nil }
func (s *checkoutServiceStub) Delete(context.Context, string) error                     { return nil }
func (s *checkoutServiceStub) CreateCheckout(_ context.Context, userID, _, _ string) (*domain.PaddleTransaction, error) {
	s.called = true
	s.userID = userID
	return &domain.PaddleTransaction{PaddleTransactionID: "txn_test", UserID: userID}, nil
}
func (s *checkoutServiceStub) GetTransaction(context.Context, string) (*domain.PaddleTransaction, error) {
	return nil, nil
}
func (s *checkoutServiceStub) HandlePaddleWebhook(context.Context, *domain.PaddleWebhookEvent) error {
	return nil
}
func (s *checkoutServiceStub) GetSubscriptionAccess(context.Context, string) (*domain.PaddleSubscriptionAccess, error) {
	return nil, nil
}

func TestCreateCheckoutRequiresAuthenticatedUser(t *testing.T) {
	service := &checkoutServiceStub{}
	handler := &Handler{PaddlePlanSrv: service, IsEnabled: true}
	request := httptest.NewRequest(http.MethodPost, "/paddle/checkout", strings.NewReader(`{"plan_id":"plan","item_id":"item"}`))
	response := httptest.NewRecorder()

	handler.CreateCheckout(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
	if service.called {
		t.Fatal("checkout service should not be called without authenticated claims")
	}
}

func TestCreateCheckoutUsesAuthenticatedUserID(t *testing.T) {
	service := &checkoutServiceStub{}
	handler := &Handler{PaddlePlanSrv: service, IsEnabled: true}
	claims := &claim.AccessTokenClaims{UserID: "authenticated-user"}
	ctx := context.WithValue(context.Background(), middlewares.AccessTokenClaimsIDKey, claims)
	request := httptest.NewRequest(http.MethodPost, "/paddle/checkout", strings.NewReader(`{"plan_id":"plan","item_id":"item","user_id":"attacker"}`)).WithContext(ctx)
	response := httptest.NewRecorder()

	handler.CreateCheckout(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; response: %s", response.Code, http.StatusCreated, response.Body.String())
	}
	if !service.called || service.userID != "authenticated-user" {
		t.Fatalf("service user ID = %q, called=%v; want authenticated-user", service.userID, service.called)
	}
}
