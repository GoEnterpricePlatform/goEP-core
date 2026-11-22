-- Catalog variations and selectable values
CREATE TABLE tb_catalog_variation (
    id         UUID PRIMARY KEY,
    name       VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_catalog_var_option (
    id           UUID PRIMARY KEY,
    variation_id UUID NOT NULL REFERENCES tb_catalog_variation(id) ON DELETE CASCADE,
    label        VARCHAR(150) NOT NULL,
    value        TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (variation_id, label)
);

-- Plans, their items, and item option selections
CREATE TABLE tb_catalog_plan (
    id          UUID PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    img_path    TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_catalog_plan_item (
    id         UUID PRIMARY KEY,
    plan_id    UUID NOT NULL REFERENCES tb_catalog_plan(id) ON DELETE CASCADE,
    img_path   TEXT,
    features   TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    status     VARCHAR(20) NOT NULL CHECK (status IN ('active', 'inactive', 'draft')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_catalog_plan_item_var_option (
    plan_item_id  UUID NOT NULL REFERENCES tb_catalog_plan_item(id) ON DELETE CASCADE,
    var_option_id UUID NOT NULL REFERENCES tb_catalog_var_option(id) ON DELETE RESTRICT,
    PRIMARY KEY (plan_item_id, var_option_id)
);

-- This model is derived from the selected variation values, but persists as
-- its own relation instead of being embedded in a plan item document.
-- ! esta por ver por que me parece que lo sacaste de aqui, OptionNoSqlModel, ese modelo creo solo era para consulta mas adelante lo veremos.
CREATE TABLE tb_catalog_plan_item_option (
    id            UUID PRIMARY KEY,
    plan_item_id  UUID NOT NULL REFERENCES tb_catalog_plan_item(id) ON DELETE CASCADE,
    name          VARCHAR(150) NOT NULL,
    var_opt_name  VARCHAR(150) NOT NULL,
    var_opt_value TEXT NOT NULL
);


-- Paddle product and price mappings owned by Catalog.
CREATE TABLE tb_catalog_paddle_plan (
    id                UUID PRIMARY KEY,
    plan_id           UUID NOT NULL UNIQUE REFERENCES tb_catalog_plan(id) ON DELETE CASCADE,
    paddle_product_id TEXT NOT NULL UNIQUE,
    "order"             INTEGER NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tb_catalog_paddle_plan_item (
    id               UUID PRIMARY KEY,
    paddle_plan_id   UUID NOT NULL REFERENCES tb_catalog_paddle_plan(id) ON DELETE CASCADE,
    plan_item_id     UUID NOT NULL UNIQUE REFERENCES tb_catalog_plan_item(id) ON DELETE CASCADE,
    paddle_price_id  TEXT NOT NULL UNIQUE,
    status           VARCHAR(20) NOT NULL CHECK (status IN ('active', 'inactive', 'draft')),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

