package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
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

	var otpCodeModel model.OtpCodeNoSqlModel
	err = r.Collection.FindOne(ctx, filter).Decode(&otpCodeModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: otp with id %s not found: %w", sharedD.ErrNotFound, id, err)
		}
		return nil, fmt.Errorf("error getting otpCode: %w", err)
	}

	var otpCode domain.OtpCode
	otpCodeModel.ToDomain(&otpCode)

	return &otpCode, nil
}
