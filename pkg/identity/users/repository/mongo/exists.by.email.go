package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ExistsByEmail checks if a user with the given email already exists in the collection.	
func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}

	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %s", err.Error())
	}

	return count > 0, nil
}
