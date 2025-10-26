package mongo

import (
	"context"
	"fmt"

	sharedD "github.com/amorindev/go-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) ConfirmEmail(ctx context.Context, userID string) error {
	oID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("%w: invalid userID :%w", sharedD.ErrIncorrectID, err)
	}

	filter := bson.M{
		"_id":            oID,
		"email_verified": false,
	}

	update := bson.M{
		"$set": bson.M{"email_verified": true},
	}

	result, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("%w: user not found", sharedD.ErrNotFound)
	}

	return nil
}
