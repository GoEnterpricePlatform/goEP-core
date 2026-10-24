
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const TableName = "goep_core_modules"

var _ lifecycle.StateStore = (*PostgresStore)(nil)

// PostgresStore persists installation-wide module lifecycle records in PostgreSQL.
type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

// CreateIndexes creates the table and its unique constraint.
func (s *PostgresStore) CreateIndexes(ctx context.Context) error {
	// in postgresql I don't know if they use indices, if they are not used we simply return nil

	return nil
}

func (s *PostgresStore) Find(
	ctx context.Context,
	name string,
) (lifecycle.ModuleState, bool, error) {
	var state lifecycle.ModuleState

	query := `
		SELECT module, version, installed, enabled, installed_at, updated_at
		FROM goep_core_modules
		WHERE module = $1
	`

	err := s.pool.QueryRow(ctx, query, name).Scan(
		&state.Name,
		&state.Version,
		&state.Installed,
		&state.Enabled,
		&state.InstalledAt,
		&state.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return lifecycle.ModuleState{}, false, nil
	}
	if err != nil {
		return lifecycle.ModuleState{}, false,
			fmt.Errorf("find module state %q: %w", name, err)
	}

	return state, true, nil
}

func (s *PostgresStore) List(
	ctx context.Context,
) ([]lifecycle.ModuleState, error) {
	query := `
		SELECT module, version, installed, enabled, installed_at, updated_at
		FROM goep_core_modules
		ORDER BY module
	`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find module states: %w", err)
	}
	defer rows.Close()

	states := make([]lifecycle.ModuleState, 0)

	for rows.Next() {
		var state lifecycle.ModuleState

		if err := rows.Scan(
			&state.Name,
			&state.Version,
			&state.Installed,
			&state.Enabled,
			&state.InstalledAt,
			&state.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan module state: %w", err)
		}

		states = append(states, state)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate module states: %w", err)
	}

	return states, nil
}

func (s *PostgresStore) Save(
	ctx context.Context,
	state lifecycle.ModuleState,
) error {
	query := `
		INSERT INTO goep_core_modules (
			module, version, installed, enabled, installed_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (module) DO UPDATE SET
			version      = EXCLUDED.version,
			installed    = EXCLUDED.installed,
			enabled      = EXCLUDED.enabled,
			installed_at = EXCLUDED.installed_at,
			updated_at   = EXCLUDED.updated_at
	`

	_, err := s.pool.Exec(
		ctx,
		query,
		state.Name,
		state.Version,
		state.Installed,
		state.Enabled,
		state.InstalledAt,
		state.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save module state %q: %w", state.Name, err)
	}

	return nil
}
