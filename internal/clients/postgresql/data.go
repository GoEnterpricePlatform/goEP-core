package postgresql

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
)

var (
	data *Data
	once sync.Once
)

// data manages the connection to database.
type Data struct {
	DB *pgx.Conn
}

// new returns a new instance of data with the database connection ready.
func New() *Data {
	//aqui se puede poner la connexion a la base de datos pero, que remos solo una ves
	once.Do(initDB)
	return data
}

func initDB() {
	db, err := getConnection()
	if err != nil {
		log.Panic(err)
	}

	// se usa para produccion?
	// eliminará las tablas
	// para entorno de produccion usar makefile
	// no causa probemas de rendimiento ?
	if os.Getenv("APP_ENV") == "dev" {
		err = MakeMigration(db)
		if err != nil {
			log.Panic(err)
		}
	}

	data = &Data{
		DB: db,
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
