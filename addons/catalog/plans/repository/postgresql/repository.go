package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/port"
	variationPort "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/port"
	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	_ port.PlanRepo                    = (*Repository)(nil)
	_ variationPort.OptionUsageChecker = (*Repository)(nil)
)

// Repository implements the catalog plan persistence contract.
type Repository struct{ Pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository { return &Repository{Pool: pool} }

func parseID(id, kind string) (uuid.UUID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid %s ID: %w", shared.ErrIncorrectID, kind, err)
	}
	return u, nil
}

func repoError(op string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return fmt.Errorf("%w: %s: %w", shared.ErrDuplicateKey, op, err)
		}
		if pgErr.Code == "23503" {
			return fmt.Errorf("%w: %s: %w", shared.ErrNotFound, op, err)
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}

func (r *Repository) planItems(ctx context.Context, planID uuid.UUID) ([]*domain.PlanItem, error) {
	rows, err := r.Pool.Query(ctx, `SELECT id::text,img_path,features,status,created_at,updated_at FROM tb_plan_item WHERE plan_id=$1 ORDER BY created_at,id`, planID)
	if err != nil {
		return nil, fmt.Errorf("find plan items: %w", err)
	}
	items := make([]*domain.PlanItem, 0)
	for rows.Next() {
		item := &domain.PlanItem{}
		if err := rows.Scan(&item.ID, &item.ImgPath, &item.Features, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan plan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("iterate plan items: %w", err)
	}
	rows.Close()
	for _, item := range items {
		itemID, _ := uuid.Parse(item.ID)
		item.Options = make([]*domain.Option, 0)
		item.VarOptionIDs = make([]string, 0)
		orows, e := r.Pool.Query(ctx, `SELECT vo.id::text,v.name,vo.label,COALESCE(vo.value,'') FROM tb_plan_item_var_option piv JOIN tb_var_option vo ON vo.id=piv.var_option_id JOIN tb_variation v ON v.id=vo.variation_id WHERE piv.plan_item_id=$1 ORDER BY v.name,vo.label`, itemID)
		if e != nil {
			return nil, fmt.Errorf("find plan item options: %w", e)
		}
		for orows.Next() {
			var oid, name, label, value string
			if e := orows.Scan(&oid, &name, &label, &value); e != nil {
				orows.Close()
				return nil, fmt.Errorf("scan plan item option: %w", e)
			}
			item.VarOptionIDs = append(item.VarOptionIDs, oid)
			item.Options = append(item.Options, &domain.Option{Name: name, VarOptName: label, VarOptValue: value})
		}
		if e := orows.Err(); e != nil {
			orows.Close()
			return nil, fmt.Errorf("iterate plan item options: %w", e)
		}
		orows.Close()
	}
	return items, nil
}
