package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/domain"
)

func (r *Repository) Exists(ctx context.Context, name domain.PermissionName) (bool, error) {
	var exists bool
	err := r.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM tb_identity_permission WHERE name = $1)`, string(name)).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check permission existence: %w", err)
	}
	return exists, nil
}