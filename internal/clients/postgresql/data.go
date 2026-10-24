package postgresql

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	data *Data
	once sync.Once
)

// data manages the connection to database.
type Data struct {
	DB   *pgx.Conn
	Pool *pgxpool.Pool
}

// new returns a new instance of data with the database connection ready.
func New(dbURI string) *Data {
	//aqui se puede poner la connexion a la base de datos pero, que remos solo una ves
	once.Do(func() {
		initDB(dbURI)
	})
	return data
}

func initDB(dbURI string) {
	db, err := getConnection(dbURI)
	if err != nil {
		log.Panic(err)
	}

	pool, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		log.Panic(err)
	}

	// is used for production or use migrations
	// Will it delete the tables?
	// for production environment use makefile
	//if os.Getenv("APP_ENV") == "dev" {
	err = MakeMigration(db)
	if err != nil {
		log.Panic(err)
	}
	//}

	data = &Data{
		DB:   db,
		Pool: pool,
	}
}

func (d *Data) Ping() error {
	ctx := context.Background()

	err := d.DB.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) Close() error {
	err := cleanDatabase(d.DB)
	if err != nil {
		return err
	}
	return d.DB.Close(context.Background())
}
