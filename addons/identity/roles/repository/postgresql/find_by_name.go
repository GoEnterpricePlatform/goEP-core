package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	role, err := r.find(ctx, `WHERE r.name=$1`, name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: role %q not found", sharedD.ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("find role by name: %w", err)
	}
	return role, nil
}