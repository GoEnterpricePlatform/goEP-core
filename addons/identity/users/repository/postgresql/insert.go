package postgresql

import (
	"context"
	"fmt"

	userDomain "github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"github.com/google/uuid"
)

func (r *Repository) Insert(ctx context.Context, u *userDomain.User) error {
	if u == nil {
		return fmt.Errorf("insert user: nil domain record")
	}
	id := uuid.New()
	if u.ID != "" {
		v, err := uuid.Parse(u.ID)
		if err != nil {
			return fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
		}
		id = v
	}
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin user insert: %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO tb_identity_user(id,email,email_verified,is_active,img_path,created_at,updated_at) VALUES($1,$2,$3,$4,$5,COALESCE($6,NOW()),COALESCE($7,NOW()))`, id, u.Email, u.EmailVerified, u.IsActive, u.ImgPath, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return userError("insert user", err)
	}
	if u.UserPassAuth != nil {
		_, err = tx.Exec(ctx, `INSERT INTO tb_identity_pass_auth(user_id,password_hash,created_at,updated_at) VALUES($1,$2,COALESCE($3,NOW()),COALESCE($4,NOW()))`, id, u.UserPassAuth.PasswordHash, u.UserPassAuth.CreatedAt, u.UserPassAuth.UpdatedAt)
		if err != nil {
			return userError("insert password auth", err)
		}
	}
	for _, raw := range u.RoleIDs {
		rid, e := uuid.Parse(raw)
		if e != nil {
			return fmt.Errorf("%w: invalid role ID %q: %w", sharedD.ErrIncorrectID, raw, e)
		}
		if _, e = tx.Exec(ctx, `INSERT INTO tb_identity_user_role(user_id,role_id) VALUES($1,$2)`, id, rid); e != nil {
			return userError("assign user role", e)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit user insert: %w", err)
	}
	u.ID = id.String()
	return nil
}