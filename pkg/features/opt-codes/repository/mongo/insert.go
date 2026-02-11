package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/features/opt-codes/domain"
	dShared "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, otp *domain.OtpCode) error {
	id := bson.NewObjectID()
	otp.ID = id

	userOID, err := bson.ObjectIDFromHex(otp.UserID.(string))
	if err != nil {
		return dShared.ErrIncorrectID
	}

	otp.UserID = userOID

	_, err = r.Collection.InsertOne(ctx, otp)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting opt code: %w", dShared.ErrDuplicateKey, err)
		}
		return fmt.Errorf("error inserting otp code: %w", err)
	}

	otp.ID = id.Hex()
	otp.UserID = userOID.Hex()

	return nil
}
