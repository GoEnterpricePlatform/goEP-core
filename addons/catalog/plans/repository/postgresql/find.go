package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) Find(ctx context.Context, rawID string) (*domain.Plan, error) {
	id, err := parseID(rawID, "plan")
	if err != nil {
		return nil, err
	}
	p := &domain.Plan{}
	err = r.Pool.QueryRow(ctx, `SELECT id::text,name,description,img_path,created_at,updated_at FROM tb_plan WHERE id=$1`, id).Scan(&p.ID, &p.Name, &p.Description, &p.ImgPath, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, sharedD.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find plan: %w", err)
	}
	p.Items, err = r.planItems(ctx, id)
	if err != nil {
		return nil, err
	}
	return p, nil
}