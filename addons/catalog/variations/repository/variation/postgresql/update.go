package postgresql

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)

func (r *Repository) Update(ctx context.Context, v *domain.Variation) error {
	id, err := parseID(v.ID, "variation")
	if err != nil {
		return err
	}
	result, err := r.Pool.Exec(ctx, `UPDATE tb_variation SET name=$2,updated_at=COALESCE($3,NOW()) WHERE id=$1`, id, v.Name, v.UpdatedAt)
	if err != nil {
		return dbError("update variation", err)
	}
	if result.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	return nil
}
