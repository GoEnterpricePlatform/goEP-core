package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/session/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/google/uuid"
)

func (r *Repository) Insert(ctx context.Context, session *domain.Session) error {
	if session == nil {
		return fmt.Errorf("insert session: nil domain record")
	}
	userID, err := uuid.Parse(session.UserID)
	if err != nil {
		return fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
	}
	id := uuid.New()
	if session.ID != "" {
		id, err = uuid.Parse(session.ID)
		if err != nil {
			return fmt.Errorf("%w: invalid session ID: %w", sharedD.ErrIncorrectID, err)
		}
	}
	_, err = r.Pool.Exec(ctx, `INSERT INTO tb_identity_session (id,user_id,refresh_token_id,refresh_token_expires_in,revoked,remember_me,created_at) VALUES ($1,$2,$3,$4,$5,$6,COALESCE($7,NOW()))`, id, userID, session.RefreshTokenID, session.RefreshTokenExpIn, session.Revoked, session.RememberMe, session.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: insert session: %w", sharedD.ErrDuplicateKey, err)
		}
		return fmt.Errorf("insert session: %w", err)
	}
	session.ID = id.String()
	return nil
}