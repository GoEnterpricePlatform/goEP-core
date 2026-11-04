package postgresql

import (
	"context"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)

func (r *Repository) Delete(ctx context.Context, rawID string) error {
	id, err := parseID(rawID, "variation")
	if err != nil {
		return err
	}
	result, err := r.Pool.Exec(ctx, `DELETE FROM tb_variation WHERE id=$1`, id)
	if err != nil {
		return dbError("delete variation", err)
	}
	if result.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	return nil
}