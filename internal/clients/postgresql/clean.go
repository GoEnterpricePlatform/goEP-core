package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// * cleand the database
func cleanDatabase(conn *pgx.Conn) error {
	return conn.Close(context.Background()); 
}