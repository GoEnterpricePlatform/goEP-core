package mongo

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

func setOptionalString(set bson.M, key string, value *string) {
	if value == nil {
		set[key] = nil
		return
	}
	set[key] = *value
}
