package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/session/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/session/model"
	dShared "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, session *domain.Session) error {
	id := bson.NewObjectID()

	userOID, err := bson.ObjectIDFromHex(session.UserID)
	if err != nil {
		return dShared.ErrIncorrectID
	}

	sessionModel := model.FromDomain(session, id, userOID)

	_, err = r.Collection.InsertOne(ctx, session)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting session: %w", dShared.ErrDuplicateKey, err)
		}
		return fmt.Errorf("error inserting session: %w", err)
	}

	sessionModel.ToDomain(session)

	return nil
}
