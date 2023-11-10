-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    provider SMALLINT NOT NULL,
    provider_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(127) NOT NULL,
    email VARCHAR(255) NOT NULL,
    profile_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system'
);

CREATE UNIQUE INDEX IF NOT EXISTS users_provider_id ON users (provider, provider_id);

CREATE INDEX IF NOT EXISTS users_username ON users (username);

CREATE INDEX IF NOT EXISTS users_email ON users (email);


CREATE TABLE IF NOT EXISTS user_statistics (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL,
    avatar_url VARCHAR(255),
    location VARCHAR(255),
    public_repositories INTEGER NOT NULL DEFAULT 0,
    followers INTEGER NOT NULL DEFAULT 0,
    following INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system',
    CONSTRAINT user_statistics_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS user_repositories (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    repository_id BIGINT NOT NULL,
    provider SMALLINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    description VARCHAR(511),
    fork BOOLEAN NOT NULL,
    fork_count BIGINT NOT NULL DEFAULT 0,
    star_count BIGINT NOT NULL DEFAULT 0,
    owner_username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_activity_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system',
    CONSTRAINT user_repositories_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS "user_repositories_repository_id" ON "user_repositories" (repository_id);

CREATE INDEX IF NOT EXISTS "user_repositories_owner_username" ON "user_repositories" (owner_username);


CREATE TABLE IF NOT EXISTS user_accesstoken (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system',
    CONSTRAINT user_accesstoken_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS user_accesstoken_user_id ON user_accesstoken (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users
DROP TABLE IF EXISTS user_statistics
DROP TABLE IF EXISTS user_repositories
DROP TABLE IF EXISTS user_accesstoken
-- +goose StatementEnd
