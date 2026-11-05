package migrations

import "embed"

// Files contains all PostgreSQL migration files.
//
//go:embed *.sql
var Files embed.FS
import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrationFiles embed.FS

func Up(ctx context.Context, db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}

	goose.SetBaseFS(migrationFiles)
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
