-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS tasks_priority_idx ON tasks USING btree (priority);
CREATE INDEX IF NOT EXISTS tasks_created_at_idx ON tasks USING btree (created_at);
CREATE INDEX IF NOT EXISTS tasks_user_id_idx ON tasks USING btree (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX tasks_user_id_idx;
DROP INDEX tasks_created_at_idx;
DROP INDEX tasks_priority_idx;
-- +goose StatementEnd
