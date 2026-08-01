package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, pProvider *domain.PaymentProvider) error {
	id := bson.NewObjectID()

	providerModel := model.FromDomain(pProvider, id)

	_, err := r.Collection.InsertOne(ctx, providerModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting payment provider: %s", sharedD.ErrDuplicateKey, err.Error())
		}
		return fmt.Errorf("error inserting payment provider: %s", err.Error())
	}

	providerModel.ToDomain(pProvider)

	return nil
}
