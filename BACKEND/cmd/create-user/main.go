package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := run(); err != nil {
		slog.Error("could not create login user", "error", err)
		os.Exit(1)
	}
}

func run() error {
	flags := flag.NewFlagSet("create-user", flag.ContinueOnError)
	name := flags.String("name", "", "user's full name")
	email := flags.String("email", "", "user's email address")
	role := flags.String("role", "clerk", "user role: clerk or administrator")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}

	*name = strings.TrimSpace(*name)
	*email = strings.TrimSpace(*email)
	*role = strings.TrimSpace(*role)
	password := os.Getenv("SEED_USER_PASSWORD")
	dsn := os.Getenv("GOOSE_DBSTRING")
	if *name == "" || *email == "" || password == "" {
		return fmt.Errorf("provide -name, -email, and SEED_USER_PASSWORD")
	}
	if dsn == "" {
		return fmt.Errorf("GOOSE_DBSTRING must point to the database to update")
	}
	if *role != "clerk" && *role != "administrator" {
		return fmt.Errorf("role must be clerk or administrator")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("check database connection: %w", err)
	}

	migrationDB := stdlib.OpenDBFromPool(pool)
	defer migrationDB.Close()
	if err := migrations.Up(ctx, migrationDB); err != nil {
		return fmt.Errorf("apply database migrations: %w", err)
	}

	var userID string
	err = pool.QueryRow(ctx, `
		INSERT INTO users (full_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (lower(email)) DO UPDATE
		SET full_name = EXCLUDED.full_name,
		    password_hash = EXCLUDED.password_hash,
		    role = EXCLUDED.role
		RETURNING id::text
	`, *name, *email, string(passwordHash), *role).Scan(&userID)
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}

	fmt.Printf("Login user ready: email=%s role=%s id=%s\n", *email, *role, userID)
	return nil
}
