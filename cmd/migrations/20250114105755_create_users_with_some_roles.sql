-- +goose Up
-- +goose StatementBegin
INSERT INTO users VALUES (1, 'admin', now(), 'admin', 'ADMIN', true);
INSERT INTO users VALUES (2, 'user', now(), 'user', 'USER', true);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = 2;
DELETE FROM users WHERE id = 1;
-- +goose StatementEnd
