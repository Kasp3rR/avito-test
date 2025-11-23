-- +goose Up
-- +goose StatementBegin

CREATE TABLE team (
    id SERIAL PRIMARY KEY,
    team_name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE users (
    user_id VARCHAR(64) PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    team_id BIGINT NOT NULL REFERENCES team (id),
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE pull_request (
    pull_request_id VARCHAR(64) PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id VARCHAR(64) NOT NULL REFERENCES users (user_id),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    merged_at TIMESTAMP
);

CREATE TABLE pull_request_reviewer (
    id SERIAL PRIMARY KEY,
    pull_request_id VARCHAR(64) NOT NULL REFERENCES pull_request (pull_request_id) ON DELETE CASCADE,
    user_id VARCHAR(64) NOT NULL REFERENCES users (user_id),
    UNIQUE (pull_request_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pull_request_reviewer;

DROP TABLE IF EXISTS pull_request;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS team;
-- +goose StatementEnd