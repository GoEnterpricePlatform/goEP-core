-- Paddle customer, subscription, and transaction records owned by Billing.
CREATE TABLE tb_billing_paddle_customer (
    id                 UUID PRIMARY KEY,
    user_id            UUID NOT NULL UNIQUE REFERENCES tb_identity_user(id) ON DELETE CASCADE,
    paddle_customer_id TEXT NOT NULL UNIQUE
);

CREATE TABLE tb_billing_paddle_subscription (
    id                       UUID PRIMARY KEY,
    paddle_subscription_id   TEXT NOT NULL UNIQUE,
    user_id                  UUID NOT NULL REFERENCES tb_identity_user(id) ON DELETE CASCADE,
    customer_id              UUID REFERENCES tb_billing_paddle_customer(id) ON DELETE SET NULL,
    plan_id                  UUID,
    plan_item_id             UUID,
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
CREATE INDEX idx_billing_paddle_subscription_user_id
    ON tb_billing_paddle_subscription(user_id);

CREATE TABLE tb_billing_paddle_transaction (
    id                       UUID PRIMARY KEY,
    paddle_transaction_id    TEXT NOT NULL UNIQUE,
    user_id                  UUID REFERENCES tb_identity_user(id) ON DELETE SET NULL,
    subscription_id          UUID REFERENCES tb_billing_paddle_subscription(id) ON DELETE SET NULL,
    customer_id              UUID REFERENCES tb_billing_paddle_customer(id) ON DELETE SET NULL,
    plan_id                  UUID NOT NULL,
    plan_item_id             UUID NOT NULL,
    status                   TEXT NOT NULL,
    currency_code            VARCHAR(3),
    subtotal                 BIGINT NOT NULL DEFAULT 0,
    tax                      BIGINT NOT NULL DEFAULT 0,
    total                    BIGINT NOT NULL DEFAULT 0,
    updated_at               TIMESTAMPTZ NOT NULL,
    last_event_at            TIMESTAMPTZ,
    payment_details_provided BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_billing_paddle_transaction_user_id
    ON tb_billing_paddle_transaction(user_id);

CREATE TABLE tb_billing_paddle_webhook_event (
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

