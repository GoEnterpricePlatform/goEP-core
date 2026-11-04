package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/domain"
	"github.com/google/uuid"
)

func (r *Repository) Insert(ctx context.Context, o *domain.VarOption) error {
	if o == nil {
		return fmt.Errorf("insert variation option: nil domain record")
	}
	id := uuid.New()
	if o.ID != "" {
		var err error
		id, err = parseID(o.ID, "variation option")
		if err != nil {
			return err
		}
	}
	variationID, err := parseID(o.VariationID, "variation")
	if err != nil {
		return err
	}
	_, err = r.Pool.Exec(ctx, `INSERT INTO tb_var_option(id,variation_id,label,value,created_at,updated_at) VALUES($1,$2,$3,$4,COALESCE($5,NOW()),COALESCE($6,NOW()))`, id, variationID, o.Label, o.Value, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		return dbError("insert variation option", err)
	}
	o.ID = id.String()
	return nil
}
