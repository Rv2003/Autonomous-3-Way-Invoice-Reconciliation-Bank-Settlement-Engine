package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/migrations"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if _, err := jwtSecret(); err != nil {
		logger.Error("authentication is not configured", "error", err)
		os.Exit(1)
	}

	corsAllowedOrigins, err := parseCORSAllowedOrigins(env.GetString("CORS_ALLOWED_ORIGINS", ""))
	if err != nil {
		logger.Error("CORS is not configured correctly", "error", err)
		os.Exit(1)
	}

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=invoice_app password=password123 dbname=invoice_reconciliation sslmode=disable"),
		},
	}
	cfg.corsAllowedOrigins = corsAllowedOrigins
	pool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		logger.Error("failed to create database connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	logger.Info("connected to database")

	migrationDB := stdlib.OpenDBFromPool(pool)
	defer migrationDB.Close()
	if err := migrations.Up(ctx, migrationDB); err != nil {
		logger.Error("failed to apply database migrations", "error", err)
		os.Exit(1)
	}

	api := application{
		config: cfg,
		db:     pool,
	}
	if err := api.run(api.mount()); err != nil {
		slog.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}
}
