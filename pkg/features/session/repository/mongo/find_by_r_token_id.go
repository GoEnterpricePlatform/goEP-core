package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/features/session/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) FindByRTokenID(ctx context.Context, rTokenID, userID string) (*domain.Session, error) {

	userOID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid userID :%w", sharedD.ErrIncorrectID, err)
	}

	filter := bson.M{
		"refresh_token_id": rTokenID,
		"user_id":          userOID,
	}

	var session domain.Session

	err = r.Collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: session not found for user %s and token %s: %w", sharedD.ErrNotFound, userID, rTokenID, err)
		}
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	if oid, ok := session.ID.(bson.ObjectID); ok {
		session.ID = oid.Hex()
	} else {
		return nil, fmt.Errorf("invalid session ID type, expected ObjectID got %T", session.ID)
	}

	session.UserID = userOID.Hex()

	return &session, nil
}
