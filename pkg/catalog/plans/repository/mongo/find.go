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

func (r *Repository) Find(ctx context.Context, id string) (*domain.Plan, error) {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, sharedD.ErrIncorrectID
	}

	var planModel model.PlanNoSqlModel
	err = r.Collection.FindOne(ctx, bson.M{"_id": oID}).Decode(&planModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("error getting plan: %w", err)
	}

	var plan domain.Plan
	planModel.ToDomain(&plan)

	return &plan, nil
}
