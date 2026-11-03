package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
)

func (r *Repository) Update(ctx context.Context, id string, role *domain.Role) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%w: invalid role ID: %w", sharedD.ErrIncorrectID, err)
	}
	tag, err := r.Pool.Exec(ctx, `UPDATE tb_identity_role SET name=$2 WHERE id=$1`, uid, role.Name)
	if err != nil {
		return roleError("update role", err)
	}
	if tag.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	return nil
}