-- MANAGED FLAVOR ONLY. Generated from managed up migrations.

-- Source: internal/modules/backend/managed/migrations/000001_module_lifecycle.up.sql
CREATE TABLE IF NOT EXISTS goep_core_modules (
    module       TEXT PRIMARY KEY,
    version      TEXT NOT NULL,
    installed    BOOLEAN NOT NULL DEFAULT FALSE,
    enabled      BOOLEAN NOT NULL DEFAULT FALSE,
    installed_at TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
