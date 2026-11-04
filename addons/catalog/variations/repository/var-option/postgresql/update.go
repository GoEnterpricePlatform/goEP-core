package postgresql

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)


func (r *Repository) Update(ctx context.Context, o *domain.VarOption) error {
	id, err := parseID(o.ID, "variation option")
	if err != nil {
		return err
	}
	variationID, err := parseID(o.VariationID, "variation")
	if err != nil {
		return err
	}
	result, err := r.Pool.Exec(ctx, `UPDATE tb_var_option SET label=$3,value=$4,updated_at=COALESCE($5,NOW()) WHERE id=$1 AND variation_id=$2`, id, variationID, o.Label, o.Value, o.UpdatedAt)
	if err != nil {
		return dbError("update variation option", err)
	}
	if result.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	return nil
}