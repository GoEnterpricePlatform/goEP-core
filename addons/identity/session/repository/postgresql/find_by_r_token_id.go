package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/session/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) FindByRTokenID(ctx context.Context, tokenID, userID string) (*domain.Session, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
	}
	var session domain.Session
	var id uuid.UUID
	var createdAt time.Time
	err = r.Pool.QueryRow(ctx, `SELECT id,user_id,refresh_token_id,refresh_token_expires_in,revoked,remember_me,created_at FROM tb_identity_session WHERE refresh_token_id=$1 AND user_id=$2`, tokenID, uid).Scan(&id, &uid, &session.RefreshTokenID, &session.RefreshTokenExpIn, &session.Revoked, &session.RememberMe, &createdAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: session not found", sharedD.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}
	session.ID, session.UserID, session.CreatedAt = id.String(), uid.String(), &createdAt
	return &session, nil
}