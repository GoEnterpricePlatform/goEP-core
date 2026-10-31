-- PostgreSQL schema for the persistent repositories currently implemented in
-- the core and addons. IDs are UUIDs in SQL; external Paddle identifiers stay
-- TEXT. Nested NoSQL models are represented by related tables where they
-- describe independent records or collections.

DROP TABLE IF EXISTS tb_paddle_webhook_event CASCADE;
DROP TABLE IF EXISTS tb_paddle_transaction CASCADE;
DROP TABLE IF EXISTS tb_paddle_subscription CASCADE;
DROP TABLE IF EXISTS tb_paddle_customer CASCADE;
DROP TABLE IF EXISTS tb_paddle_plan_item CASCADE;
DROP TABLE IF EXISTS tb_paddle_plan CASCADE;
DROP TABLE IF EXISTS tb_payment_provider_lemonsqueezy CASCADE;
DROP TABLE IF EXISTS tb_payment_provider_stripe CASCADE;
DROP TABLE IF EXISTS tb_payment_provider CASCADE;
DROP TABLE IF EXISTS tb_plan_item_var_option CASCADE;
DROP TABLE IF EXISTS tb_plan_item_option CASCADE;
DROP TABLE IF EXISTS tb_plan_item CASCADE;
DROP TABLE IF EXISTS tb_plan CASCADE;
DROP TABLE IF EXISTS tb_var_option CASCADE;
DROP TABLE IF EXISTS tb_variation CASCADE;
DROP TABLE IF EXISTS tb_post CASCADE;
DROP TABLE IF EXISTS tb_identity_otp_code CASCADE;
DROP TABLE IF EXISTS tb_identity_session CASCADE;
DROP TABLE IF EXISTS tb_identity_user_role CASCADE;
DROP TABLE IF EXISTS tb_identity_role_permission CASCADE;
DROP TABLE IF EXISTS tb_identity_role CASCADE;
DROP TABLE IF EXISTS tb_identity_permission CASCADE;
DROP TABLE IF EXISTS tb_identity_pass_auth CASCADE;
DROP TABLE IF EXISTS tb_identity_user CASCADE;
DROP TABLE IF EXISTS goep_core_modules CASCADE;

CREATE TABLE goep_core_modules (
    module       TEXT PRIMARY KEY,
    version      TEXT NOT NULL,
    installed    BOOLEAN NOT NULL DEFAULT FALSE,
    enabled      BOOLEAN NOT NULL DEFAULT FALSE,
    installed_at TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Identity
CREATE TABLE tb_identity_user (
    id             UUID PRIMARY KEY,
    email          VARCHAR(250) NOT NULL UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    img_path       VARCHAR(300),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Password authentication is an optional one-to-one user credential.
CREATE TABLE tb_identity_pass_auth (
    user_id       UUID PRIMARY KEY REFERENCES tb_identity_user(id) ON DELETE CASCADE,
    password_hash VARCHAR(256) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_identity_permission (
    id   UUID PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('SYSTEM', 'VIEW'))
);

CREATE TABLE tb_identity_role (
    id   UUID PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE
);

CREATE TABLE tb_identity_role_permission (
    role_id       UUID NOT NULL REFERENCES tb_identity_role(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES tb_identity_permission(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE tb_identity_user_role (
    user_id UUID NOT NULL REFERENCES tb_identity_user(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES tb_identity_role(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE tb_identity_session (
    id                   UUID PRIMARY KEY,
    user_id              UUID NOT NULL REFERENCES tb_identity_user(id) ON DELETE CASCADE,
    refresh_token_id     TEXT NOT NULL UNIQUE,
    refresh_token_expires_in BIGINT NOT NULL,
    revoked              BOOLEAN NOT NULL DEFAULT FALSE,
    remember_me          BOOLEAN NOT NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_identity_otp_code (
    id         UUID PRIMARY KEY,
    user_id    UUID NOT NULL REFERENCES tb_identity_user(id) ON DELETE CASCADE,
    otp_code   VARCHAR(10) NOT NULL,
    purpose    VARCHAR(40) NOT NULL,
    used       BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_identity_otp_code_user_id ON tb_identity_otp_code(user_id);

-- Posts
CREATE TABLE tb_post (
    id         UUID PRIMARY KEY,
    title      TEXT NOT NULL,
    "desc"     TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Catalog variations and selectable values
CREATE TABLE tb_variation (
    id         UUID PRIMARY KEY,
    name       VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_var_option (
    id           UUID PRIMARY KEY,
    variation_id UUID NOT NULL REFERENCES tb_variation(id) ON DELETE CASCADE,
    label        VARCHAR(150) NOT NULL,
    value        TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (variation_id, label)
);

-- Plans, their items, and item option selections
CREATE TABLE tb_plan (
    id          UUID PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    img_path    TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_plan_item (
    id         UUID PRIMARY KEY,
    plan_id    UUID NOT NULL REFERENCES tb_plan(id) ON DELETE CASCADE,
    img_path   TEXT,
    features   TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    status     VARCHAR(20) NOT NULL CHECK (status IN ('active', 'inactive', 'draft')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_plan_item_var_option (
    plan_item_id  UUID NOT NULL REFERENCES tb_plan_item(id) ON DELETE CASCADE,
    var_option_id UUID NOT NULL REFERENCES tb_var_option(id) ON DELETE RESTRICT,
    PRIMARY KEY (plan_item_id, var_option_id)
);

-- This model is derived from the selected variation values, but persists as
-- its own relation instead of being embedded in a plan item document.
-- ! esta por ver por que me parece que lo sacaste de aqui, OptionNoSqlModel, ese modelo creo solo era para consulta mas adelante lo veremos.
CREATE TABLE tb_plan_item_option (
    id            UUID PRIMARY KEY,
    plan_item_id  UUID NOT NULL REFERENCES tb_plan_item(id) ON DELETE CASCADE,
    name          VARCHAR(150) NOT NULL,
    var_opt_name  VARCHAR(150) NOT NULL,
    var_opt_value TEXT NOT NULL
);

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

-- Paddle catalog mapping, customers, subscriptions, transactions and events.
CREATE TABLE tb_paddle_plan (
    id                UUID PRIMARY KEY,
    plan_id           UUID NOT NULL UNIQUE REFERENCES tb_plan(id) ON DELETE CASCADE,
    paddle_product_id TEXT NOT NULL UNIQUE,
    "order"             INTEGER NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_paddle_plan_item (
    id               UUID PRIMARY KEY,
    paddle_plan_id   UUID NOT NULL REFERENCES tb_paddle_plan(id) ON DELETE CASCADE,
    plan_item_id     UUID NOT NULL UNIQUE REFERENCES tb_plan_item(id) ON DELETE CASCADE,
    paddle_price_id  TEXT NOT NULL UNIQUE,
    status           VARCHAR(20) NOT NULL CHECK (status IN ('active', 'inactive', 'draft')),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_paddle_customer (
    id                 UUID PRIMARY KEY,
    user_id            UUID NOT NULL UNIQUE REFERENCES tb_identity_user(id) ON DELETE CASCADE,
    paddle_customer_id TEXT NOT NULL UNIQUE
);

-- ! Here why the three do not use not null in the reference id and set null, of course check the subscription model why it is added as "omitempty", check that the omitempty is used correctly.
CREATE TABLE tb_paddle_subscription (
    id                       UUID PRIMARY KEY,
    paddle_subscription_id   TEXT NOT NULL UNIQUE,
    user_id                  UUID NOT NULL REFERENCES tb_identity_user(id) ON DELETE CASCADE,
    customer_id              UUID REFERENCES tb_paddle_customer(id) ON DELETE SET NULL,
    plan_id                  UUID REFERENCES tb_plan(id) ON DELETE SET NULL,
    plan_item_id             UUID REFERENCES tb_plan_item(id) ON DELETE SET NULL,
    status                   TEXT NOT NULL,
    current_period_starts_at TIMESTAMPTZ,
    current_period_ends_at   TIMESTAMPTZ,
    next_billed_at           TIMESTAMPTZ,
    canceled_at              TIMESTAMPTZ,
    scheduled_change_action  TEXT,
    scheduled_change_effective_at TIMESTAMPTZ,
    updated_at               TIMESTAMPTZ NOT NULL,
    last_event_at            TIMESTAMPTZ
);
CREATE INDEX idx_paddle_subscription_user_id ON tb_paddle_subscription(user_id);

CREATE TABLE tb_paddle_transaction (
    id                       UUID PRIMARY KEY,
    paddle_transaction_id    TEXT NOT NULL UNIQUE,
    user_id                  UUID REFERENCES tb_identity_user(id) ON DELETE SET NULL,
    subscription_id          UUID REFERENCES tb_paddle_subscription(id) ON DELETE SET NULL,
    customer_id              UUID REFERENCES tb_paddle_customer(id) ON DELETE SET NULL,
    plan_id                  UUID NOT NULL REFERENCES tb_plan(id) ON DELETE RESTRICT,
    plan_item_id             UUID NOT NULL REFERENCES tb_plan_item(id) ON DELETE RESTRICT,
    status                   TEXT NOT NULL,
    currency_code            VARCHAR(3),
    subtotal                 BIGINT NOT NULL DEFAULT 0,
    tax                      BIGINT NOT NULL DEFAULT 0,
    total                    BIGINT NOT NULL DEFAULT 0,
    updated_at               TIMESTAMPTZ NOT NULL,
    last_event_at            TIMESTAMPTZ,
    payment_details_provided BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_paddle_transaction_user_id ON tb_paddle_transaction(user_id);

CREATE TABLE tb_paddle_webhook_event (
    id                    UUID PRIMARY KEY,
    event_id              TEXT NOT NULL UNIQUE,
    event_type            TEXT NOT NULL,
    resource_id           TEXT,
    occurred_at           TIMESTAMPTZ,
    status                TEXT NOT NULL,
    received_at           TIMESTAMPTZ NOT NULL,
    last_received_at      TIMESTAMPTZ NOT NULL,
    processing_started_at TIMESTAMPTZ NOT NULL,
    processed_at          TIMESTAMPTZ,
    failed_at             TIMESTAMPTZ,
    attempts              INTEGER NOT NULL DEFAULT 0,
    last_error            TEXT
);
