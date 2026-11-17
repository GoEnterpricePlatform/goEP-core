-- Posts
CREATE TABLE IF NOT EXISTS tb_post (
    id         UUID PRIMARY KEY,
    title      TEXT NOT NULL,
    "desc"     TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

