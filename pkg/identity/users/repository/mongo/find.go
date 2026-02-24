package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Find(ctx context.Context, id string) (*domain.User, error) {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid userID :%w", sharedD.ErrIncorrectID, err)
	}

	var userModel model.UserNoSqlModel
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oID}}).Decode(&userModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: user with id %s not found: %w", sharedD.ErrNotFound, id, err)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	var user domain.User
	userModel.ToDomain(&user)

	return &user, nil
}
