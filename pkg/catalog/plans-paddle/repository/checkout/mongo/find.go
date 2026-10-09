package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *CheckoutRepository) Find(ctx context.Context, transactionID string) (*domain.PaddleCheckout, error) {
	var checkoutModel model.PaddleCheckoutModel
	if err := r.Collection.FindOne(ctx, bson.M{"_id": transactionID}).Decode(&checkoutModel); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("error finding Paddle checkout: %w", err)
	}
	checkout := &domain.PaddleCheckout{}
	checkoutModel.ToDomain(checkout)
	return checkout, nil
}
