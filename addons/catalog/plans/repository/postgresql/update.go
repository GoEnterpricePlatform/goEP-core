package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
)

func (r *Repository) Update(ctx context.Context, rawID string, p *domain.Plan) error {
	id, err := parseID(rawID, "plan")
	if err != nil {
		return err
	}
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin plan update: %w", err)
	}
	defer tx.Rollback(ctx)
	res, err := tx.Exec(ctx, `UPDATE tb_plan SET name=$2,description=$3,img_path=$4,updated_at=COALESCE($5,NOW()) WHERE id=$1`, id, p.Name, p.Description, p.ImgPath, p.UpdatedAt)
	if err != nil {
		return repoError("update plan", err)
	}
	if res.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	// The service sends the complete editable item set; remove rows omitted from it.
	keep := make(map[uuid.UUID]struct{}, len(p.Items))
	for _, item := range p.Items {
		if item == nil {
			continue
		}
		itemID, e := parseID(item.ID, "plan item")
		if e != nil {
			return e
		}
		keep[itemID] = struct{}{}
	}
	rows, err := tx.Query(ctx, `SELECT id FROM tb_plan_item WHERE plan_id=$1`, id)
	if err != nil {
		return fmt.Errorf("list plan items before update: %w", err)
	}
	var removed []uuid.UUID
	for rows.Next() {
		var itemID uuid.UUID
		if err := rows.Scan(&itemID); err != nil {
			rows.Close()
			return fmt.Errorf("scan plan item before update: %w", err)
		}
		if _, exists := keep[itemID]; !exists {
			removed = append(removed, itemID)
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate plan items before update: %w", err)
	}
	rows.Close()
	for _, itemID := range removed {
		if _, err := tx.Exec(ctx, `DELETE FROM tb_plan_item WHERE id=$1 AND plan_id=$2`, itemID, id); err != nil {
			return repoError("remove plan item", err)
		}
	}
	for _, item := range p.Items {
		if item == nil {
			continue
		}
		itemID, e := parseID(item.ID, "plan item")
		if e != nil {
			return e
		}
		features := item.Features
		if features == nil {
			features = []string{}
		}
		res, e := tx.Exec(ctx, `UPDATE tb_plan_item SET img_path=$3,features=$4,updated_at=COALESCE($5,NOW()) WHERE id=$1 AND plan_id=$2`, itemID, id, item.ImgPath, features, item.UpdatedAt)
		if e != nil {
			return repoError("update plan item", e)
		}
		if res.RowsAffected() == 0 {
			return sharedD.ErrNotFound
		}
		if _, e = tx.Exec(ctx, `DELETE FROM tb_plan_item_var_option WHERE plan_item_id=$1`, itemID); e != nil {
			return fmt.Errorf("replace plan item options: %w", e)
		}
		for _, raw := range item.VarOptionIDs {
			oid, e := parseID(raw, "variation option")
			if e != nil {
				return e
			}
			if _, e = tx.Exec(ctx, `INSERT INTO tb_plan_item_var_option(plan_item_id,var_option_id) VALUES($1,$2)`, itemID, oid); e != nil {
				return repoError("assign plan item option", e)
			}
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit plan update: %w", err)
	}
	return nil
}