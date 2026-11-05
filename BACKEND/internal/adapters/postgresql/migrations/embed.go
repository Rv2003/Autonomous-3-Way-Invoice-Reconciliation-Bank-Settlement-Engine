package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

// Files contains all PostgreSQL migration files.
//
//go:embed *.sql
var Files embed.FS

func Up(ctx context.Context, db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}

	goose.SetBaseFS(Files)
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
