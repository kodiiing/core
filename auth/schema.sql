-- Dialect: PostgreSQL

-- "users" table contains information that are used to define the user.
-- This data must not change often.
CREATE TABLE IF NOT EXISTS "users" (
    id BIGSERIAL PRIMARY KEY,
    provider TINYINT NOT NULL,
    provider_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(127) NOT NULL,
    email VARCHAR(255) NOT NULL,
    profile_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system'
);

CREATE UNIQUE INDEX IF NOT EXISTS "users_provider_id" ON "users" (provider, provider_id);

CREATE INDEX IF NOT EXISTS "users_username" ON "users" (username);

CREATE INDEX IF NOT EXISTS "users_email" ON "users" (email);

-- "user_statistics" table contains information from the
-- provider that is used as an extra metadata for the user.
-- Information such as location, avatar url, public repository,
-- followers and following belongs here.
CREATE TABLE IF NOT EXISTS "user_statistics" (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL,
    avatar_url VARCHAR(255),
    location VARCHAR(255),
    public_repositories INTEGER NOT NULL DEFAULT 0,
    followers INTEGER NOT NULL DEFAULT 0,
    following INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system',
    CONSTRAINT "user_statistics_user_id_fkey" FOREIGN KEY (user_id) REFERENCES "users" (id) ON DELETE CASCADE
);

-- "user_repositories" table contains information about the repositories
-- that the user owns or even had contributed to.
CREATE TABLE IF NOT EXISTS "user_repositories" (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    repository_id BIGINT NOT NULL,
    provider TINYINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    description VARCHAR(511),
    fork BOOLEAN NOT NULL,
    fork_count BIGINT NOT NULL DEFAULT 0,
    star_count BIGINT NOT NULL DEFAULT 0,
    owner_username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_activity_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system',
    CONSTRAINT "user_repositories_user_id_fkey" FOREIGN KEY (user_id) REFERENCES "users" (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS "user_repositories_repository_id" ON "user_repositories" (repository_id);

CREATE INDEX IF NOT EXISTS "user_repositories_owner_username" ON "user_repositories" (owner_username);

CREATE TABLE IF NOT EXISTS "user_accesstoken" (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system',
    CONSTRAINT "user_accesstoken_user_id_fkey" FOREIGN KEY (user_id) REFERENCES "users" (id) ON DELETE CASCADE
);
