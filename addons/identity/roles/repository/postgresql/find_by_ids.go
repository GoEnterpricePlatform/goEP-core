package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindByIDs(ctx context.Context, ids []string) ([]*domain.Role, error) {
	if len(ids) == 0 {
		return []*domain.Role{}, nil
	}
	uids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		u, e := uuid.Parse(id)
		if e != nil {
			return nil, fmt.Errorf("%w: invalid role ID %q: %w", sharedD.ErrIncorrectID, id, e)
		}
		uids = append(uids, u)
	}
	return r.findMany(ctx, `WHERE r.id=ANY($1::uuid[])`, uids)
}