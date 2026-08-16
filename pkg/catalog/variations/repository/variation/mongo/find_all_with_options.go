package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) FindAllWithOptions(ctx context.Context) ([]*domain.Variation, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "var_options"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "variation_id"},
			{Key: "as", Value: "options"},
		}}},
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error getting variations with options: %w", err)
	}
	defer cursor.Close(ctx)

	var variations []*domain.Variation
	for cursor.Next(ctx) {
		var variationModel model.VariationNoSqlModel
		err := cursor.Decode(&variationModel)
		if err != nil {
			return nil, fmt.Errorf("error decoding variation: %w", err)
		}
		var variation domain.Variation
		variationModel.ToDomain(&variation)

		variations = append(variations, &variation)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cursor: %w", err)
	}

	return variations, nil
}
