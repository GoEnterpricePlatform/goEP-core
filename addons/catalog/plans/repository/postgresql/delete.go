package postgresql

import (
	"context"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)


func (r *Repository) Delete(ctx context.Context, rawID string) error {
	id, err := parseID(rawID, "plan")
	if err != nil {
		return err
	}
	res, err := r.Pool.Exec(ctx, `DELETE FROM tb_plan WHERE id=$1`, id)
	if err != nil {
		return repoError("delete plan", err)
	}
	if res.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	return nil
}