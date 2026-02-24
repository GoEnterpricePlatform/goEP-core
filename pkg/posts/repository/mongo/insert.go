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

func (r *Repository) Insert(ctx context.Context, posts *domain.Post) error {
	id := bson.NewObjectID()

	postModel := model.FromDomain(posts, id)

	_, err := r.Collection.InsertOne(ctx, postModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting post: %s", sharedD.ErrDuplicateKey, err.Error())
		}
		return fmt.Errorf("error inserting post: %s", err.Error())
	}

	postModel.ToDomain(posts)

	return nil
}
