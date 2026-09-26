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

	pipeline := mongo.Pipeline{
		// Match plan by ID.
		{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: oID},
		}}},

		// Join with var_options.
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "catalog_var-options"},
			{Key: "localField", Value: "items.var_option_ids"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "all_var_options"},
		}}},

		// Join with variations.
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "catalog_variations"},
			{Key: "localField", Value: "all_var_options.variation_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "all_variations"},
		}}},

		// Rebuild plan items with their options.
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
													{Key: "$in", Value: bson.A{
														"$$opt._id",
														"$$item.var_option_ids",
													}},
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
																	{Key: "$eq", Value: bson.A{
																		"$$v._id",
																		"$$opt.variation_id",
																	}},
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
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate pipeline: %w", err)
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("failed to read plan from cursor: %w", err)
		}

		return nil, sharedD.ErrNotFound
	}

	var planModel model.PlanNoSqlModel

	if err := cursor.Decode(&planModel); err != nil {
		return nil, fmt.Errorf("failed to decode plan from cursor: %w", err)
	}

	var plan domain.Plan
	planModel.ToDomain(&plan)

	return &plan, nil
}
