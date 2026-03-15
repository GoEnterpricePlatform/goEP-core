package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) Search(ctx context.Context, query string, limit int) ([]*domain.Post, error) {
	filter := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error searching post: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []*domain.Post
	for cursor.Next(ctx) {
		var postModel model.PostNoSqlModel
		err := cursor.Decode(&postModel)
		if err != nil {
			return nil, fmt.Errorf("error decoding post: %w", err)
		}
		var post domain.Post
		postModel.ToDomain(&post)

		posts = append(posts, &post)
	}

	return posts, nil
}
