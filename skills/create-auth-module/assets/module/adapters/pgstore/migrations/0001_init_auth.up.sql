BEGIN;

CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users
(
    user_uuid uuid NOT NULL,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    password_hash varchar(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_uuid),
    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE TABLE auth.refresh_tokens
(
    token_uuid uuid NOT NULL,
    user_uuid uuid NOT NULL REFERENCES auth.users (user_uuid) ON DELETE CASCADE,
    token_hash varchar(255) NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    revoked_at timestamptz,
    PRIMARY KEY (token_uuid),
    CONSTRAINT refresh_tokens_hash_unique UNIQUE (token_hash)
);

-- Refresh, logout, and reuse detection revoke by user.
CREATE INDEX refresh_tokens_user_uuid_idx ON auth.refresh_tokens (user_uuid);

COMMIT;
