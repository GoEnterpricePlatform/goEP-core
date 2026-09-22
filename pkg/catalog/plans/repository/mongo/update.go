package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) Update(ctx context.Context, id string, plan *domain.Plan) error {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return sharedD.ErrIncorrectID
	}
	
	items := make([]bson.D, 0, len(plan.Items))
	
	for _, item := range plan.Items {
		if item == nil {
			continue
		}
		
		itemID, err := bson.ObjectIDFromHex(item.ID)
		if err != nil {
			return sharedD.ErrIncorrectID
		}

		varOptionIDs := make([]bson.ObjectID, 0, len(item.VarOptionIDs))

		for _, varOptionID := range item.VarOptionIDs {
			varOptionOID, err := bson.ObjectIDFromHex(varOptionID)
			if err != nil {
				return sharedD.ErrIncorrectID
			}

			varOptionIDs = append(varOptionIDs, varOptionOID)
		}

		items = append(items, bson.D{
			{Key: "_id", Value: itemID},
			{Key: "img_path", Value: item.ImgPath},
			{Key: "features", Value: item.Features},
			{Key: "status", Value: item.Status},
			{Key: "var_option_ids", Value: varOptionIDs},
			{Key: "created_at", Value: item.CreatedAt},
			{Key: "updated_at", Value: item.UpdatedAt},
		})
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: plan.Name},
				{Key: "description", Value: plan.Description},
				{Key: "img_path", Value: plan.ImgPath},
				{Key: "items", Value: items},
				{Key: "updated_at", Value: plan.UpdatedAt},
			},
		},
	}

	filter := bson.D{
		{Key: "_id", Value: oID},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedPlan model.PlanNoSqlModel

	err = r.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedPlan)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return sharedD.ErrNotFound
		}

		return fmt.Errorf("failed to update plan: %w", err)
	}

	updatedPlan.ToDomain(plan)

	return nil
}
