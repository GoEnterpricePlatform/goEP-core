package postgresql

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func getConnection() (client *pgx.Conn, err error) {
	// set time out
	dbURI := os.Getenv("PG_DB_URI")
	if dbURI == "" {
		return nil, errors.New("PG_DB_URI not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	//connect to the PostgreSql sercer
	conn, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
