package postgresql

import (
	"context"
	"errors"
	"fmt"

	userDomain "github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Find(ctx context.Context, id string) (*userDomain.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
	}
	u, err := r.find(ctx, `WHERE u.id=$1`, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, sharedD.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	return u, nil
}