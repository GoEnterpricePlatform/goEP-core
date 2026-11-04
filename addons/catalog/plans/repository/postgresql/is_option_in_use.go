package postgresql

import "context"

// mongo podías eliminar un registro que esta referenciado en otro lugar se puede hacer dos cosas
// manejarlo desde postgresql me parece que si tiene esa opcion y en IsOptionInUse retornar true y nil
// o de ambas formas
func (r *Repository) IsOptionInUse(ctx context.Context, rawID string) (bool, error) {
	id, err := parseID(rawID, "variation option")
	if err != nil {
		return false, err
	}
	var used bool
	err = r.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tb_plan_item_var_option WHERE var_option_id=$1)`, id).Scan(&used)
	return used, err
}