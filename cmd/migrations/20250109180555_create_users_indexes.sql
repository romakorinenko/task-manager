-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS users_login_idx ON users USING btree (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX users_login_idx;
-- +goose StatementEnd
