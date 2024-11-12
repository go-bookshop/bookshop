BEGIN;

CREATE TABLE IF NOT EXISTS tokens
(
    id      BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hash    BYTEA                       NOT NULL UNIQUE,
    user_id BIGINT                      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expiry  TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope   TEXT                        NOT NULL
);

CREATE INDEX idx_tokens_user_id ON tokens (user_id);
CREATE INDEX idx_tokens_hash ON tokens (hash);

COMMIT;