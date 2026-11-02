package migrations

import "embed"

// Files contains all PostgreSQL migration files.
//
//go:embed *.sql
var Files embed.FS
