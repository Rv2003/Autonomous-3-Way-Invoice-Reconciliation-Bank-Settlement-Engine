package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/migrations"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbString := os.Getenv("GOOSE_DBSTRING")
	if dbString == "" {
		log.Fatal("GOOSE_DBSTRING is not set")
	}

	db, err := sql.Open("pgx", dbString)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	goose.SetBaseFS(migrations.Files)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set Goose dialect: %v", err)
	}

	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("failed to apply database migrations: %v", err)
	}

	log.Println("database migrations completed successfully")
}
