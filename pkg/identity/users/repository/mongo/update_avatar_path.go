package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) UpdateAvatarPath(ctx context.Context, userID, imgPath string, updatedAt time.Time) error {
	oID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("%w: failed to convert userID to ObjectID :%w", domain.ErrIncorrectID, err)
	}

	filter := bson.M{"_id": oID}
	update := bson.M{"$set": bson.M{"img_path": imgPath, "updated_at": updatedAt}}

	result, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user avatar path in database: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}
