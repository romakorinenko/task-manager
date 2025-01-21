-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE users_sequence start 100;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE users_sequence;
-- +goose StatementEnd
