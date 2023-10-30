-- +goose Up
CREATE TABLE
    IF NOT EXISTS authors(
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        access_token TEXT CONSTRAINT unique_access_token UNIQUE NOT NULL,
        profile_url VARCHAR(255) NOT NULL,
        picture_url VARCHAR(255) NOT NULL,
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp
    )
CREATE TABLE
    IF NOT EXISTS comments (
        id BIGSERIAL PRIMARY KEY,
        content TEXT NOT NULL,
        author_id BIGSERIAL NOT NULL,
        created_at timestamp default current_timestamp,
        CONSTRAINT fk_author_comments foreign key (author_id) REFERENCES authors(id)
    )
CREATE TABLE
    IF NOT EXISTS hacks(
        id BIGSERIAL PRIMARY KEY,
        author_id BIGSERIAL NOT NULL,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        tags VARCHAR [] NOT NULL,
        upvotes BIGINT null,
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp,
        CONSTRAINT fk_author_comments foreign key (author_id) REFERENCES authors(id)
    )
CREATE TABLE
    IF NOT EXISTS hack_comments (
        id BIGSERIAL PRIMARY KEY,
        hack_id BIGSERIAL NOT NULL,
        comment_id BIGSERIAL NOT NULL,
        parent_id BIGINT null,
        created_at timestamp default current_timestamp,
        CONSTRAINT fk_hack foreign key (hack_id) REFERENCES hacks(id),
        CONSTRAINT fk_comment foreign key (comment_id) REFERENCES comments(id)
    )

-- +goose Down
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS hacks;
DROP TABLE IF EXISTS hack_comments;

