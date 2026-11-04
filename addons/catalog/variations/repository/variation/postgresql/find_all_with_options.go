package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindAllWithOptions(ctx context.Context) ([]*domain.Variation, error) {
	rows, err := r.Pool.Query(ctx, `SELECT id::text,name,created_at,updated_at FROM tb_variation ORDER BY name,id`)
	if err != nil {
		return nil, fmt.Errorf("find variations: %w", err)
	}
	variations := make([]*domain.Variation, 0)
	for rows.Next() {
		v := &domain.Variation{}
		if err := rows.Scan(&v.ID, &v.Name, &v.CreatedAt, &v.UpdatedAt); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan variation: %w", err)
		}
		variations = append(variations, v)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("iterate variations: %w", err)
	}
	rows.Close()
	for _, v := range variations {
		id, _ := uuid.Parse(v.ID)
		optionRows, err := r.Pool.Query(ctx, `SELECT id::text,variation_id::text,label,value,created_at,updated_at FROM tb_var_option WHERE variation_id=$1 ORDER BY label,id`, id)
		if err != nil {
			return nil, fmt.Errorf("find variation options: %w", err)
		}
		v.Options = make([]*domain.VarOption, 0)
		for optionRows.Next() {
			o := &domain.VarOption{}
			if err := optionRows.Scan(&o.ID, &o.VariationID, &o.Label, &o.Value, &o.CreatedAt, &o.UpdatedAt); err != nil {
				optionRows.Close()
				return nil, fmt.Errorf("scan variation option: %w", err)
			}
			v.Options = append(v.Options, o)
		}
		if err := optionRows.Err(); err != nil {
			optionRows.Close()
			return nil, fmt.Errorf("iterate variation options: %w", err)
		}
		optionRows.Close()
	}
	return variations, nil
}
