package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/domain"
	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindByIDs(ctx context.Context, ids []string) ([]*domain.Permission, error) {
	if len(ids) == 0 {
		return []*domain.Permission{}, nil
	}
	uids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		parsed, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid permission ID %q: %w", shared.ErrIncorrectID, id, err)
		}
		uids = append(uids, parsed)
	}
	rows, err := r.Pool.Query(ctx, `SELECT id::text, name, type FROM tb_identity_permission WHERE id = ANY($1::uuid[])`, uids)
	if err != nil {
		return nil, fmt.Errorf("find permissions by IDs: %w", err)
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
