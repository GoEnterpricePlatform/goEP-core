package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/features/users/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Find(ctx context.Context, id string) (*domain.User, error) {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid userID :%w", sharedD.ErrIncorrectID, err)
	}

	var user domain.User
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oID}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: user with id %s not found: %w", sharedD.ErrNotFound, id, err)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	user.ID = oID.Hex()

	return &user, nil
}
