package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"

	"github.com/google/uuid"
)

func (r *Repository) Insert(ctx context.Context, role *domain.Role) error {
	if role == nil {
		return fmt.Errorf("insert role: nil domain record")
	}
	id := uuid.New()
	if role.ID != "" {
		parsed, err := uuid.Parse(role.ID)
		if err != nil {
			return fmt.Errorf("%w: invalid role ID: %w", sharedD.ErrIncorrectID, err)
		}
		id = parsed
	}
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin role insert: %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO tb_identity_role (id,name) VALUES ($1,$2)`, id, role.Name)
	if err != nil {
		return roleError("insert role", err)
	}
	for _, raw := range role.PermissionIDs {
		pid, e := uuid.Parse(raw)
		if e != nil {
			return fmt.Errorf("%w: invalid permission ID %q: %w", sharedD.ErrIncorrectID, raw, e)
		}
		if _, e = tx.Exec(ctx, `INSERT INTO tb_identity_role_permission (role_id,permission_id) VALUES ($1,$2)`, id, pid); e != nil {
			return roleError("assign role permission", e)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit role insert: %w", err)
	}
	role.ID = id.String()
	return nil
}