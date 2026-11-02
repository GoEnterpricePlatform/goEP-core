package postgresql

import (
	"context"
	"fmt"
	"time"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"

	"github.com/google/uuid"
)

func (r *Repository) UpdateAvatarPath(ctx context.Context, id, path string, updated time.Time) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
	}
	tag, err := r.Pool.Exec(ctx, `UPDATE tb_identity_user SET img_path=$2,updated_at=$3 WHERE id=$1`, uid, path, updated)
	if err != nil {
		return fmt.Errorf("update user avatar path: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	return nil
}