package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model/customer"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *CustomerRepository) FindByUserID(ctx context.Context, userID string) (*domain.PaddleCustomer, error) {
	var doc model.PaddleCustomerModel
	if err := r.Collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("error finding Paddle customer by user ID: %w", err)
	}
	return doc.ToDomain(), nil
}
