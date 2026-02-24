package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/model"
	dShared "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, otp *domain.OtpCode) error {
	id := bson.NewObjectID()

	optCodeModel, err := model.FromDomain(otp, id)
	if err != nil {
		return fmt.Errorf("error mapping otpCode to mongo model: %w", err)
	}

	_, err = r.Collection.InsertOne(ctx, optCodeModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting opt code: %w", dShared.ErrDuplicateKey, err)
		}
		return fmt.Errorf("error inserting otp code: %w", err)
	}

	optCodeModel.ToDomain(otp)

	return nil
}
