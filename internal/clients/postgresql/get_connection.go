package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func getConnection(dbURI string) (client *pgx.Conn, err error) {
	// set time out
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	//connect to the PostgreSql sercer
	conn, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
