package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, paddlePlan *domain.PaddlePlan) error {
	oID := bson.NewObjectID()

	paddlePlanModel, err := model.FromDomainPaddlePlan(paddlePlan, oID)
	if err != nil {
		return fmt.Errorf("error mapping paddle plan to mongo model: %w", err)
	}

	_, err = r.Collection.InsertOne(ctx, paddlePlanModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf(
				"%w: error inserting paddle plan: %w",
				sharedD.ErrDuplicateKey,
				err,
			)
		}

		return fmt.Errorf("error inserting paddle plan: %w", err)
	}

	paddlePlanModel.ToDomain(paddlePlan)

	return nil
}
