package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) FindByIDs(ctx context.Context, ids []string) ([]*domain.VarOption, error) {
	objectIDs := make([]bson.ObjectID, 0, len(ids))
	for _, idStr := range ids {
		oid, err := bson.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid id %s: %w", sharedD.ErrIncorrectID, idStr, err)
		}
		objectIDs = append(objectIDs, oid)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find VarOptions by IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*domain.VarOption
	for cursor.Next(ctx) {
		optionModel := &model.VarOptionNoSqlModel{}

		if err := cursor.Decode(optionModel); err != nil {
			return nil, fmt.Errorf("failed to decode VarOption: %w", err)
		}

		option := &domain.VarOption{}
		optionModel.ToDomain(option)

		results = append(results, option)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate VarOptions: %w", err)
	}
	return results, nil
}
