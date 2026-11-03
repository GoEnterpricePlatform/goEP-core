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

