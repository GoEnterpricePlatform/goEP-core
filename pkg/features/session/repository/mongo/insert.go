package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-tmpl/pkg/features/session/domain"
	dShared "github.com/amorindev/go-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, session *domain.Session) error {
	id := bson.NewObjectID()
	session.ID = id

	userOID, err := bson.ObjectIDFromHex(session.UserID.(string))
	if err != nil {
		return dShared.ErrIncorrectID
	}

	session.UserID = userOID

	_, err = r.Collection.InsertOne(ctx, session)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting session: %w", dShared.ErrDuplicateKey, err)
		}
		return fmt.Errorf("error inserting session: %w", err)
	}

	session.ID = id.Hex()
	session.UserID = userOID.Hex()

	return nil
}
