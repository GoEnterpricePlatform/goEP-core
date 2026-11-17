package main

import (
	"log"
	"os"

	"github.com/GoEnterpricePlatform/goEP-core/internal/migrations"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to load .env file in migration main: %v", err)
	}

	pgDbUri := os.Getenv("PG_DB_URI")
	if pgDbUri == "" {
		log.Fatalf("Environment variable %s is required but not set", "PG_DB_URI")
	}

	if len(os.Args) != 2 {
		log.Fatal("usage: go run ./cmd/migration <up | down>")
	}

	switch os.Args[1] {
	case "up":
		if err := migrations.Up(pgDbUri); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := migrations.Down(pgDbUri); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
