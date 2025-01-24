-- +goose Up
-- +goose StatementBegin
CREATE table IF NOT EXISTS users
(
    id         BIGINT PRIMARY KEY,
    login      VARCHAR(255) NOT NULL UNIQUE,
    created_at DATE DEFAULT CURRENT_DATE,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(255) NOT NULL,
    active     BOOLEAN DEFAULT TRUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
