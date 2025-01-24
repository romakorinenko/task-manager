-- +goose Up
-- +goose StatementBegin
CREATE table IF NOT EXISTS tasks
(
    id          BIGINT PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    priority    BIGINT       NOT NULL,
    status      VARCHAR(255) NOT NULL,
    created_at  DATE DEFAULT CURRENT_DATE,
    updated_at  DATE DEFAULT CURRENT_DATE,
    user_id     BIGINT REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
