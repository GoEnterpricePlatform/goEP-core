package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) FindAll(ctx context.Context, limit int64, page int64) ([]*domain.Plan, error) {
	skip := (page - 1) * limit

	pipeline := mongo.Pipeline{
		// join with var_options
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "catalog_var-options"},
			{Key: "localField", Value: "items.var_option_ids"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "all_var_options"},
		}}},
		// join with variations
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "catalog_variations"},
			{Key: "localField", Value: "all_var_options.variation_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "all_variations"},
		}}},
		// rebuilding plan_items with their options
		{{Key: "$addFields", Value: bson.D{
			{Key: "items", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$items"},
					{Key: "as", Value: "item"},
					{Key: "in", Value: bson.D{
						{Key: "$mergeObjects", Value: bson.A{
							"$$item",
							bson.D{
								{Key: "options", Value: bson.D{
									{Key: "$map", Value: bson.D{
										{Key: "input", Value: bson.D{
											{Key: "$filter", Value: bson.D{
												{Key: "input", Value: "$all_var_options"},
												{Key: "as", Value: "opt"},
												{Key: "cond", Value: bson.D{
													{Key: "$in", Value: bson.A{"$$opt._id", "$$item.var_option_ids"}},
												}},
											}},
										}},
										{Key: "as", Value: "opt"},
										{Key: "in", Value: bson.D{
											{Key: "name", Value: bson.D{
												{Key: "$first", Value: bson.D{
													{Key: "$map", Value: bson.D{
														{Key: "input", Value: bson.D{
															{Key: "$filter", Value: bson.D{
																{Key: "input", Value: "$all_variations"},
																{Key: "as", Value: "v"},
																{Key: "cond", Value: bson.D{
																	{Key: "$eq", Value: bson.A{"$$v._id", "$$opt.variation_id"}},
																}},
															}},
														}},
														{Key: "as", Value: "v"},
														{Key: "in", Value: "$$v.name"},
													}},
												}},
											}},
											{Key: "var_opt_name", Value: "$$opt.label"},
											{Key: "var_opt_value", Value: "$$opt.value"},
										}},
									}},
								}},
							},
						}},
					}},
				}},
			}},
		}}},

		// pagination
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	var planModels []*model.PlanNoSqlModel

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate pipeline: %w", err)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &planModels); err != nil {
		return nil, fmt.Errorf("failed to decode plans from cursor: %w", err)
	}

	// Convert the NoSQL models to domain objects explicitly.
	plans := make([]*domain.Plan, 0, len(planModels))

	for _, planModel := range planModels {
		if planModel == nil {
			continue
		}

		plan := &domain.Plan{}
		planModel.ToDomain(plan)

		plans = append(plans, plan)
	}

	return plans, nil
}
