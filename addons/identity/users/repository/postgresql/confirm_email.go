package postgresql

import (
	"context"
	"fmt"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"

	"github.com/google/uuid"
)

func (r *Repository) ConfirmEmail(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
	}
	tag, err := r.Pool.Exec(ctx, `UPDATE tb_identity_user SET email_verified=TRUE,updated_at=NOW() WHERE id=$1`, uid)
	if err != nil {
		return fmt.Errorf("confirm user email: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	return nil
}