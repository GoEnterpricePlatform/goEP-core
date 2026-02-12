package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/users/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/users/model"
	sharedDomain "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// * Insert adds a new user, returning ErrDuplicateKey if the user already exists.
func (r *Repository) Insert(ctx context.Context, user *domain.User) error {
	id := bson.NewObjectID()

	userModel, err := model.FromDomain(user, id)
	if err != nil {
		return fmt.Errorf("error mapping user to mongo model: %w", err)
	}

	_, err = r.Collection.InsertOne(ctx, userModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting user: %s", sharedDomain.ErrDuplicateKey, err.Error())
		}
		return fmt.Errorf("error inserting user: %s", err.Error())
	}

	userModel.ToDomain(user)

	return nil
}
