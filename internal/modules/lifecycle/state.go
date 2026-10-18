package lifecycle

import (
	"context"
	"errors"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
)

var (
	ErrModuleNotInstalled = errors.New("module is not installed")
	ErrModuleNotAvailable = errors.New("module is not available in this build")
)

// ModuleState stores installation-wide lifecycle state for one module.
type ModuleState struct {
	Name        string    `bson:"module" json:"module"`
	Version     string    `bson:"version" json:"version"`
	Installed   bool      `bson:"installed" json:"installed"`
	Enabled     bool      `bson:"enabled" json:"enabled"`
	InstalledAt time.Time `bson:"installed_at" json:"installed_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// StateStore is the persistence contract used by the module lifecycle manager.
type StateStore interface {
	Find(ctx context.Context, name string) (ModuleState, bool, error)
	List(ctx context.Context) ([]ModuleState, error)
	Save(ctx context.Context, state ModuleState) error
}

// ModuleStatus combines discovered metadata with persisted installation state.
type ModuleStatus struct {
	Manifest  modules.Manifest
	Installed bool
	Enabled   bool
}
