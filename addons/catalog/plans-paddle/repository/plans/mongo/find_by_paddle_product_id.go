package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/model/plan"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) FindByPaddleProductID(ctx context.Context, paddleProductID string) (*domain.PaddlePlan, error) {
	var paddlePlanModel model.PaddlePlanNoSqlModel
	err := r.Collection.FindOne(ctx, bson.M{"paddle_product_id": paddleProductID}).Decode(&paddlePlanModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("find Paddle plan by product ID: %w", err)
	}
	paddlePlan := &domain.PaddlePlan{}
	paddlePlanModel.ToDomain(paddlePlan)
	return paddlePlan, nil
}
