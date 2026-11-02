package postgresql

import (
	"context"
	"fmt"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"

	"github.com/google/uuid"
)

func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return false, fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
	}
	var ok bool
	err = r.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tb_identity_user WHERE id=$1)`, uid).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("check user existence: %w", err)
	}
	return ok, nil
}