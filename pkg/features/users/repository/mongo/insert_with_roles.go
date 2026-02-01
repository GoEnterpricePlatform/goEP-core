package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-tmpl/pkg/features/users/domain"
	sharedD "github.com/amorindev/go-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) InsertWithRoles(ctx context.Context, user *domain.User) error {
	id := bson.NewObjectID()
	user.ID = id

	var rolesIDsInsert []interface{}
	for _, idStr := range user.RoleIDs {
		oID, err := bson.ObjectIDFromHex(idStr.(string))
		if err != nil {
			return fmt.Errorf("%w: invalid role id %s :%w", sharedD.ErrIncorrectID, idStr, err)
		}
		rolesIDsInsert = append(rolesIDsInsert, oID)
	}

	user.RoleIDs = rolesIDsInsert

	_, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting user with roles: %s", sharedD.ErrDuplicateKey, err.Error())
		}
		return fmt.Errorf("error inserting user with roles: %s", err.Error())
	}

	user.ID = id.Hex()

	return nil
}
