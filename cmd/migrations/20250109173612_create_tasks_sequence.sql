-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE tasks_sequence start 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE tasks_sequence;
-- +goose StatementEnd
