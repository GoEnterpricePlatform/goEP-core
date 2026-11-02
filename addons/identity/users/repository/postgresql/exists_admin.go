package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
)

func (r *Repository) ExistsAdmin(ctx context.Context) (bool, error) {
	var ok bool
	err := r.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tb_identity_user_role ur JOIN tb_identity_role r ON r.id=ur.role_id WHERE r.name=$1)`, string(domain.RoleSystemAdmin)).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("check admin existence: %w", err)
	}
	return ok, nil
}