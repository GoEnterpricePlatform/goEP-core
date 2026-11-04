package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/domain"
	"github.com/jackc/pgx/v5"
)


func (r *Repository) FindByIDs(ctx context.Context, rawIDs []string) ([]*domain.VarOption, error) {
	options := make([]*domain.VarOption, 0, len(rawIDs))
	for _, rawID := range rawIDs {
		id, err := parseID(rawID, "variation option")
		if err != nil {
			return nil, err
		}
		o := &domain.VarOption{}
		err = r.Pool.QueryRow(ctx, `SELECT id::text,variation_id::text,label,value,created_at,updated_at FROM tb_var_option WHERE id=$1`, id).Scan(&o.ID, &o.VariationID, &o.Label, &o.Value, &o.CreatedAt, &o.UpdatedAt)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("find variation option: %w", err)
		}
		options = append(options, o)
	}
	return options, nil
}