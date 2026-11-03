package postgresql

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
)

func MakeMigration(conn *pgx.Conn) error {

	path := filepath.Join(
		"internal",
		"clients",
		"postgresql",
		"model",
		"schema.sql",
	)

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = conn.Exec(ctx, string(b))
	if err != nil {
		return err
	}

	return nil
}
