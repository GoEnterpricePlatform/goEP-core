package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/domain"
)

func (r *Repository) FindByNames(ctx context.Context, names []domain.PermissionName) ([]*domain.Permission, error) {
	if len(names) == 0 {
		return []*domain.Permission{}, nil
	}
	values := make([]string, 0, len(names))
	for _, name := range names {
		values = append(values, string(name))
	}
	rows, err := r.Pool.Query(ctx, `SELECT id::text, name, type FROM tb_identity_permission WHERE name = ANY($1::text[])`, values)
	if err != nil {
		return nil, fmt.Errorf("find permissions by names: %w", err)
	}
	defer rows.Close()
	
	permissions := make([]*domain.Permission, 0)
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Type); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		permissions = append(permissions, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate permissions: %w", err)
	}
	return permissions, nil
}