package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
	"go.mongodb.org/mongo-driver/v2/bson"
	driver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const MongoCollectionName = "goep_core_modules"

var _ lifecycle.StateStore = (*MongoStore)(nil)

// MongoStore persists installation-wide module lifecycle records in MongoDB.
type MongoStore struct {
	collection *driver.Collection
}

func NewMongoStore(database *driver.Database) *MongoStore {
	return &MongoStore{collection: database.Collection(MongoCollectionName)}
}

// CreateIndexes ensures one state document per module.
func (s *MongoStore) CreateIndexes(ctx context.Context) error {
	_, err := s.collection.Indexes().CreateOne(ctx, driver.IndexModel{
		Keys:    bson.D{{Key: "module", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_module_state_name"),
	})
	if err != nil {
		return fmt.Errorf("create module state index: %w", err)
	}
	return nil
}

func (s *MongoStore) Find(ctx context.Context, name string) (lifecycle.ModuleState, bool, error) {
	var state lifecycle.ModuleState
	err := s.collection.FindOne(ctx, bson.D{{Key: "module", Value: name}}).Decode(&state)
	if err == driver.ErrNoDocuments {
		return lifecycle.ModuleState{}, false, nil
	}
	if err != nil {
		return lifecycle.ModuleState{}, false, fmt.Errorf("find module state %q: %w", name, err)
	}
	return state, true, nil
}

func (s *MongoStore) List(ctx context.Context) ([]lifecycle.ModuleState, error) {
	cursor, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("find module states: %w", err)
	}
	defer cursor.Close(ctx)

	states := make([]lifecycle.ModuleState, 0)
	if err := cursor.All(ctx, &states); err != nil {
		return nil, fmt.Errorf("decode module states: %w", err)
	}
	return states, nil
}

func (s *MongoStore) Save(ctx context.Context, state lifecycle.ModuleState) error {
	filter := bson.D{{Key: "module", Value: state.Name}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "module", Value: state.Name},
		{Key: "version", Value: state.Version},
		{Key: "installed", Value: state.Installed},
		{Key: "enabled", Value: state.Enabled},
		{Key: "installed_at", Value: state.InstalledAt},
		{Key: "updated_at", Value: state.UpdatedAt},
	}}}
	if _, err := s.collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true)); err != nil {
		return fmt.Errorf("save module state %q: %w", state.Name, err)
	}
	return nil
}
