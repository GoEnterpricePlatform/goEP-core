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

	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Fatal("usage: go run ./cmd/migration <up | down> [lite | managed]")
	}
	flavor := "lite"
	if len(os.Args) == 3 {
		flavor = os.Args[2]
	}
	if flavor != "lite" && flavor != "managed" {
		log.Fatalf("unknown migration flavor %q: supported flavors are lite and managed", flavor)
	}

	switch os.Args[1] {
	case "up":
		var err error
		if flavor == "managed" {
			err = migrations.UpManaged(pgDbUri)
		} else {
			err = migrations.Up(pgDbUri)
		}
		if err != nil {
			log.Fatal(err)
		}
	case "down":
		var err error
		if flavor == "managed" {
			err = migrations.DownManaged(pgDbUri)
		} else {
			err = migrations.Down(pgDbUri)
		}
		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
