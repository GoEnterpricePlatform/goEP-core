package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/domain"
	"github.com/google/uuid"
)

func (r *Repository) Insert(ctx context.Context, v *domain.Variation) error {
	if v == nil {
		return fmt.Errorf("insert variation: nil domain record")
	}
	id := uuid.New()
	if v.ID != "" {
		var err error
		id, err = parseID(v.ID, "variation")
		if err != nil {
			return err
		}
	}
	_, err := r.Pool.Exec(ctx, `INSERT INTO tb_variation(id,name,created_at,updated_at) VALUES($1,$2,COALESCE($3,NOW()),COALESCE($4,NOW()))`, id, v.Name, v.CreatedAt, v.UpdatedAt)
	if err != nil {
		return dbError("insert variation", err)
	}
	v.ID = id.String()
	return nil
}
