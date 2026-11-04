# MongoDB Repository Notes

## Indexes

In MongoDB we don't need to create indexes to prevent duplicates, so `CreateIndexes` is a no-op:

```go
func (r *Repository) CreateIndexes() error {
	return nil
}
```

## Deleting Referenced Records

MongoDB lets you delete a record that is still referenced somewhere else (there are no foreign key constraints). There are two ways to handle this:

1. Handle it from PostgreSQL, which I believe supports this option (foreign key constraints), and have `IsOptionInUse` return `true, nil`.
2. Do both: enforce it in PostgreSQL and also check it in `IsOptionInUse`.

```go
func (r *Repository) IsOptionInUse(ctx context.Context, rawID string) (bool, error) {
	id, err := parseID(rawID, "variation option")
	if err != nil {
		return false, err
	}
	var used bool
	err = r.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tb_plan_item_var_option WHERE var_option_id=$1)`, id).Scan(&used)
	return used, err
}
```