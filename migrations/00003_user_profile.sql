-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS user_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    join_reason SMALLINT NOT NULL DEFAULT 4,
    join_reason_other TEXT NULL,
    coded_before BOOLEAN NOT NULL DEFAULT FALSE,
    languages TEXT NULL,
    target TEXT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system'
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_join_reason ON user_profiles (join_reason);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_profiles_join_reason;
DROP TABLE IF EXISTS user_profiles;

-- +goose StatementEnd