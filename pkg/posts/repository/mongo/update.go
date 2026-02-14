package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/posts/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) Update(ctx context.Context, id string, post *domain.Post) error {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return sharedD.ErrIncorrectID
	}

	update := bson.M{
		"$set": bson.M{
			"title":      post.Title,
			"desc":       post.Desc,
			"updated_at": post.UpdatedAt,
		},
	}

	filter := bson.M{"_id": oID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = r.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return sharedD.ErrNotFound
		}
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}
