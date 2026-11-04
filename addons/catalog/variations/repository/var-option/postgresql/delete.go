package postgresql

import (
	"context"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)

func (r *Repository) Delete(ctx context.Context, rawID, rawVariationID string) error {
	id, err := parseID(rawID, "variation option")
	if err != nil {
		return err
	}
	variationID, err := parseID(rawVariationID, "variation")
	if err != nil {
		return err
	}
	result, err := r.Pool.Exec(ctx, `DELETE FROM tb_var_option WHERE id=$1 AND variation_id=$2`, id, variationID)
	if err != nil {
		return dbError("delete variation option", err)
	}
	if result.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	return nil
}
