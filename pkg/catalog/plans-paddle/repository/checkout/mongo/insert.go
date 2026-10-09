package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
)

func (r *CheckoutRepository) Insert(ctx context.Context, checkout *domain.PaddleCheckout) error {
	checkoutModel := model.FromDomainPaddleCheckout(checkout)
	if _, err := r.Collection.InsertOne(ctx, checkoutModel); err != nil {
		return fmt.Errorf("error saving Paddle checkout: %w", err)
	}
	return nil
}
