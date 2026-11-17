package core

import "embed"

// FS contains the PostgreSQL migrations owned by GoEP core.
//
//go:embed *.sql
var FS embed.FS
