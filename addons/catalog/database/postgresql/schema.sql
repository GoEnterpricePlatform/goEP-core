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
