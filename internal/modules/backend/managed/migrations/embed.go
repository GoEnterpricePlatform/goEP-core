package migrations

import "embed"

// FS contains PostgreSQL migrations used only by the managed flavor.
//
//go:embed *.sql
var FS embed.FS
