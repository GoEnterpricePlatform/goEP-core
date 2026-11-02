package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) FindAll(ctx context.Context) ([]*domain.Role, error) {
	return r.findMany(ctx, "", nil)
}
func (r *Repository) findMany(ctx context.Context, where string, arg any) ([]*domain.Role, error) {
	query := `SELECT r.id::text,r.name,COALESCE(array_agg(rp.permission_id::text) FILTER (WHERE rp.permission_id IS NOT NULL),ARRAY[]::text[]) FROM tb_identity_role r LEFT JOIN tb_identity_role_permission rp ON rp.role_id=r.id ` + where + ` GROUP BY r.id,r.name ORDER BY r.name`
	var rows pgx.Rows
	var err error
	if where == "" {
		rows, err = r.Pool.Query(ctx, query)
	} else {
		rows, err = r.Pool.Query(ctx, query, arg)
	}
	if err != nil {
		return nil, fmt.Errorf("find roles: %w", err)
	}
	defer rows.Close()
	roles := make([]*domain.Role, 0)
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.PermissionIDs); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, &role)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate roles: %w", err)
	}
	return roles, nil
}