package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, plan *domain.Plan) error {
	oID := bson.NewObjectID()
	
	planModel, err := model.FromDomainPlan(plan, oID)
	if err != nil {
		return fmt.Errorf("error mapping plan to mongo model: %w", err)
	}
	
	_, err = r.Collection.InsertOne(ctx, planModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting plan: %w", sharedD.ErrDuplicateKey, err)
		}
		return fmt.Errorf("error inserting plan: %w", err)
	}
	
	planModel.ToDomain(plan)

	return nil
}