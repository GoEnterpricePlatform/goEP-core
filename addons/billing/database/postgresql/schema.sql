-- Billing provider configuration, with one-to-one provider-specific settings.
CREATE TABLE tb_payment_provider (
    id         UUID PRIMARY KEY,
    name       VARCHAR(30) NOT NULL UNIQUE CHECK (name IN ('stripe', 'lemonsqueezy')),
    enabled    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_payment_provider_stripe (
    provider_id                  UUID PRIMARY KEY REFERENCES tb_payment_provider(id) ON DELETE CASCADE,
    publishable_key              TEXT NOT NULL,
    secret_key_encrypted         TEXT NOT NULL,
    webhook_secret_encrypted     TEXT NOT NULL
);

CREATE TABLE tb_payment_provider_lemonsqueezy (
    provider_id                  UUID PRIMARY KEY REFERENCES tb_payment_provider(id) ON DELETE CASCADE,
    api_key_encrypted            TEXT NOT NULL,
    store_id                     TEXT NOT NULL,
    webhook_secret_encrypted     TEXT NOT NULL
);

