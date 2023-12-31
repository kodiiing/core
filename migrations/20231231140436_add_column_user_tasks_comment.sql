-- +goose Up
-- +goose StatementBegin

ALTER TABLE user_tasks ADD COLUMN IF NOT EXISTS comments NOT NULL DEFAULT '';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE user_tasks DROP COULUMN IF EXISTS comments;

-- +goose StatementEnd
