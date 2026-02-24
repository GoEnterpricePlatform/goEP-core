package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/model"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) FindAll(ctx context.Context) ([]*domain.Post, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error getting posts: %w", err)
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
