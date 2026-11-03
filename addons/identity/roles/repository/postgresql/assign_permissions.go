package postgresql

import (
	"context"
	"errors"
	"fmt"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
)

func (r *Repository) AssignPermissions(ctx context.Context, name string, ids []string) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin permission assignment: %w", err)
	}

	defer tx.Rollback(ctx)

	var roleID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM tb_identity_role WHERE name=$1`, name).Scan(&roleID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: role %q not found", sharedD.ErrNotFound, name)
	}
	if err != nil {
		return fmt.Errorf("find role for permission assignment: %w", err)
	}
	for _, raw := range ids {
		pid, e := uuid.Parse(raw)
		if e != nil {
			return fmt.Errorf("%w: invalid permission ID %q: %w", sharedD.ErrIncorrectID, raw, e)
		}
		_, e = tx.Exec(ctx, `INSERT INTO tb_identity_role_permission(role_id,permission_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, roleID, pid)
		if e != nil {
			return fmt.Errorf("assign permission: %w", e)
		}
	}
	return tx.Commit(ctx)
}