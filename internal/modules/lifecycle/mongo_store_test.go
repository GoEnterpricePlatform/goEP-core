package lifecycle

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	driver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestMongoStorePersistsStateAcrossClients(t *testing.T) {
	uri := os.Getenv("GOEP_TEST_MONGODB_URI")
	if uri == "" {
		t.Skip("set GOEP_TEST_MONGODB_URI to run the MongoDB persistence integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	databaseName := "goep_core_lifecycle_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")

	firstClient, err := driver.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("connect first MongoDB client: %v", err)
	}
	firstDatabase := firstClient.Database(databaseName)
	store := NewMongoStore(firstDatabase)
	if err := store.CreateIndexes(ctx); err != nil {
		_ = firstClient.Disconnect(context.Background())
		t.Fatalf("CreateIndexes() error = %v", err)
	}
	want := ModuleState{
		Name:        "catalog",
		Version:     "1.2.3",
		Installed:   true,
		Enabled:     true,
		InstalledAt: time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC),
		UpdatedAt:   time.Date(2026, time.February, 3, 4, 5, 6, 0, time.UTC),
	}
	if err := store.Save(ctx, want); err != nil {
		_ = firstClient.Disconnect(context.Background())
		t.Fatalf("Save() error = %v", err)
	}
	if err := firstClient.Disconnect(ctx); err != nil {
		t.Fatalf("disconnect first MongoDB client: %v", err)
	}

	secondClient, err := driver.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("connect second MongoDB client: %v", err)
	}
	defer secondClient.Disconnect(context.Background())
	secondDatabase := secondClient.Database(databaseName)
	defer secondDatabase.Drop(context.Background())
	store = NewMongoStore(secondDatabase)

	got, found, err := store.Find(ctx, want.Name)
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if !found || got != want {
		t.Fatalf("Find() = %#v, found=%t; want %#v", got, found, want)
	}

	states, err := store.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(states) != 1 || states[0] != want {
		t.Fatalf("List() = %#v, want one state %#v", states, want)
	}
}
