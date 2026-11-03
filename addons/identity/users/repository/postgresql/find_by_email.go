package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	authDomain "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/domain"
	userDomain "github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) FindByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	u, err := r.find(ctx, `WHERE u.email=$1`, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: user with email %s not found", sharedD.ErrNotFound, email)
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return u, nil
}

func (r *Repository) find(ctx context.Context, where string, arg any) (*userDomain.User, error) {
	var u userDomain.User
	var hash *string
	var passCreated, passUpdated *time.Time
	err := r.Pool.QueryRow(ctx, `SELECT u.id::text,u.email,u.email_verified,u.is_active,u.img_path,u.created_at,u.updated_at,p.password_hash,p.created_at,p.updated_at FROM tb_identity_user u LEFT JOIN tb_identity_pass_auth p ON p.user_id=u.id `+where, arg).Scan(&u.ID, &u.Email, &u.EmailVerified, &u.IsActive, &u.ImgPath, &u.CreatedAt, &u.UpdatedAt, &hash, &passCreated, &passUpdated)
	if err != nil {
		return nil, err
	}
	if hash != nil {
		u.UserPassAuth = &authDomain.UserPasswordAuth{PasswordHash: *hash, CreatedAt: passCreated, UpdatedAt: passUpdated}
	}
	rows, err := r.Pool.Query(ctx, `SELECT role_id::text FROM tb_identity_user_role WHERE user_id=$1`, u.ID)
	if err != nil {
		return nil, fmt.Errorf("find user roles: %w", err)
	}
	defer rows.Close()
	u.RoleIDs = make([]string, 0)
	for rows.Next() {
		var rid string
		if err := rows.Scan(&rid); err != nil {
			return nil, fmt.Errorf("scan user role: %w", err)
		}
		u.RoleIDs = append(u.RoleIDs, rid)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user roles: %w", err)
	}
	return &u, nil
}