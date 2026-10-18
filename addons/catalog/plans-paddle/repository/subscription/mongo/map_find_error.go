package mongo

import (
	"errors"
	"fmt"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func mapFindError(err error, message string) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return sharedD.ErrNotFound
	}
	return fmt.Errorf("%s: %w", message, err)
}
