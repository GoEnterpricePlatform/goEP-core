package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Get(ctx context.Context, id string) (*domain.Post, error) {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, sharedD.ErrIncorrectID
	}

	var postModel model.PostNoSqlModel
	err = r.Collection.FindOne(ctx, bson.M{"_id": oID}).Decode(&postModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("error getting post: %w", err)
	}

	var post domain.Post
	postModel.ToDomain(&post)

	return &post, nil
}
