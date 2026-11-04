package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	"github.com/google/uuid"
)

func (r *Repository) Insert(ctx context.Context, p *domain.Plan) error {
	if p == nil {
		return fmt.Errorf("insert plan: nil domain record")
	}
	id := uuid.New()
	if p.ID != "" {
		var err error
		id, err = parseID(p.ID, "plan")
		if err != nil {
			return err
		}
	}
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin plan insert: %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO tb_plan(id,name,description,img_path,created_at,updated_at) VALUES($1,$2,$3,$4,COALESCE($5,NOW()),COALESCE($6,NOW()))`, id, p.Name, p.Description, p.ImgPath, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return repoError("insert plan", err)
	}
	p.ID = id.String()
	for _, item := range p.Items {
		if item == nil {
			continue
		}
		features := item.Features
		if features == nil {
			features = []string{}
		}
		itemID := uuid.New()
		if item.ID != "" {
			itemID, err = parseID(item.ID, "plan item")
			if err != nil {
				return err
			}
		}
		_, err = tx.Exec(ctx, `INSERT INTO tb_plan_item(id,plan_id,img_path,features,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,COALESCE($6,NOW()),COALESCE($7,NOW()))`, itemID, id, item.ImgPath, features, item.Status, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return repoError("insert plan item", err)
		}
		item.ID = itemID.String()
		for _, raw := range item.VarOptionIDs {
			optionID, e := parseID(raw, "variation option")
			if e != nil {
				return e
			}
			if _, e = tx.Exec(ctx, `INSERT INTO tb_plan_item_var_option(plan_item_id,var_option_id) VALUES($1,$2)`, itemID, optionID); e != nil {
				return repoError("assign plan item option", e)
			}
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit plan insert: %w", err)
	}
	return nil
}