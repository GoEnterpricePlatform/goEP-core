package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) FindAll(ctx context.Context, limit, page int64) ([]*domain.PaddlePlan, error) {
	opts := options.Find().SetSkip((page - 1) * limit).SetLimit(limit).SetSort(bson.D{{Key: "order", Value: 1}, {Key: "created_at", Value: -1}, {Key: "_id", Value: -1}})
	cursor, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("error listing paddle plans: %w", err)
	}
	defer cursor.Close(ctx)
	var models []*model.PaddlePlanNoSqlModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, fmt.Errorf("error decoding paddle plans: %w", err)
	}
	plans := make([]*domain.PaddlePlan, 0, len(models))
	for _, m := range models {
		if m == nil {
			continue
		}
		plan := &domain.PaddlePlan{}
		m.ToDomain(plan)
		plans = append(plans, plan)
	}
	return plans, nil
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, fmt.Errorf("error counting paddle plans: %w", err)
	}
	return count, nil
}
