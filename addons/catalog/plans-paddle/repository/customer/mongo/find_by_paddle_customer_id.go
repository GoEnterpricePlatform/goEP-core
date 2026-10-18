package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/model/customer"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *CustomerRepository) FindByPaddleCustomerID(ctx context.Context, paddleCustomerID string) (*domain.PaddleCustomer, error) {
	var doc model.PaddleCustomerModel
	if err := r.Collection.FindOne(ctx, bson.M{"paddle_customer_id": paddleCustomerID}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("error finding Paddle customer by Paddle customer ID: %w", err)
	}
	return doc.ToDomain(), nil
}
