package postgresql

import (
	"context"
	"fmt"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"


	"github.com/google/uuid"
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%w: invalid role ID: %w", sharedD.ErrIncorrectID, err)
	}
	tag, err := r.Pool.Exec(ctx, `DELETE FROM tb_identity_role WHERE id=$1`, uid)
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	return nil
}