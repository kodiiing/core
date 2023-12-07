-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(511) NOT NULL,
    difficulty SMALLINT NOT NULL,
    content TEXT NOT NULL,
    author BIGINT REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(63) NOT NULL DEFAULT 'system',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(63) NOT NULL DEFAULT 'system'
);

CREATE INDEX IF NOT EXISTS idx_tasks_difficulty ON tasks (difficulty);


CREATE TABLE IF NOT EXISTS user_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT REFERENCES tasks(id),
    user_id BIGINT REFERENCES users(id),
    status SMALLINT NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ NULL,
    finished_at TIMESTAMPTZ NULL,
    satisfaction_level SMALLINT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_task_id ON user_tasks (task_id, user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tasks_difficulty;
DROP TABLE IF EXISTS tasks;

DROP INDEX IF EXISTS idx_user_task_id;
DROP TABLE IF EXISTS user_tasks;
-- +goose StatementEnd
