package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func getConnection(dbURI string) (*mongo.Client, error) {
	timeDuration := time.Second * 10

	// Configure client options with connection timeout.
	cp := options.ClientOptions{
		ConnectTimeout: &timeDuration,
	}

	// Connect to MongoDB using the configured options and URI.
	c, err := mongo.Connect(&cp, options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, err
	}
	return c, nil
}
