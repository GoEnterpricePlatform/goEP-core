package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) FindIDs(ctx context.Context, names []string) ([]string, error) {
	if len(names) == 0 {
		return []string{}, nil
	}

	filter := bson.M{"name": bson.M{"$in": names}}

	opts := options.Find().SetProjection(bson.M{"_id": 1})

	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles by names: %w", err)
	}
	defer cursor.Close(ctx)

	var ids []string
	for cursor.Next(ctx) {
		var result struct {
			ID bson.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode role document: %w", err)
		}
		ids = append(ids, result.ID.Hex())
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error while iterating role documents: %w", err)
	}

	return ids, nil
}
