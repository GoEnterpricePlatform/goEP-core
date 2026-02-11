package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/features/opt-codes/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Find(ctx context.Context, id, userID string) (*domain.OtpCode, error) {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid otpID :%w", sharedD.ErrIncorrectID, err)
	}

	userOID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid userID :%w", sharedD.ErrIncorrectID, err)
	}

	filter := bson.D{
		{Key: "_id", Value: oID},
		{Key: "user_id", Value: userOID},
	}

	var otp domain.OtpCode
	err = r.Collection.FindOne(ctx, filter).Decode(&otp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: otp with id %s not found: %w", sharedD.ErrNotFound, id, err)
		}
		return nil, fmt.Errorf("error getting otpCode: %w", err)
	}

	otp.ID = oID.Hex()
	otp.UserID = userOID.Hex()

	return &otp, nil
}
