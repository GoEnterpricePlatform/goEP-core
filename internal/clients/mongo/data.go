package mongo

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Package mongo provides a singleton pattern for MongoDB client initialization and lifecycle management.

// Data holds the MongoDB client instance.
var (
	data *Data
	once sync.Once
)

// Data holds the MongoDB client instance.
type Data struct {
	DB *mongo.Client
}

// New initializes the MongoDB client once and returns the singleton instance.
func New(dbURI string) *Data {
	once.Do(func() {
		initDB(dbURI)
	})
	return data
}

// initDB connects to MongoDB and sets the client in the singleton Data struct.
func initDB(dbURI string) {
	db, err := getConnection(dbURI)
	if err != nil {
		log.Fatal(err)
	}

	data = &Data{
		DB: db,
	}
}

// Ping checks the connection to MongoDB.
func (data *Data) Ping() error {
	return data.DB.Ping(context.Background(), nil)
}

// Close disconnects the MongoDB client.
func (data *Data) Close() error {
	return data.DB.Disconnect(context.Background())
}
