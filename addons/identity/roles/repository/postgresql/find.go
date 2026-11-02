package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Find(ctx context.Context, id string) (*domain.Role, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid role ID: %w", sharedD.ErrIncorrectID, err)
	}
	role, err := r.find(ctx, `WHERE r.id=$1`, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, sharedD.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find role: %w", err)
	}
	return role, nil
}

func (r *Repository) find(ctx context.Context, where string, arg any) (*domain.Role, error) {
	var role domain.Role
	err := r.Pool.QueryRow(ctx, `SELECT r.id::text,r.name,COALESCE(array_agg(rp.permission_id::text) FILTER (WHERE rp.permission_id IS NOT NULL),ARRAY[]::text[]) FROM tb_identity_role r LEFT JOIN tb_identity_role_permission rp ON rp.role_id=r.id `+where+` GROUP BY r.id,r.name`, arg).Scan(&role.ID, &role.Name, &role.PermissionIDs)
	return &role, err
}